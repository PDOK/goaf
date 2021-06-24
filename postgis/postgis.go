package postgis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"oaf-server/provider"
	"time"

	"github.com/go-spatial/geom/encoding/geojson"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type IdNotFoundError struct {
	err string
}

type Postgis struct {
	ApplicationId string
	UserVersion   string
	db            *sqlx.DB
	Collections   []provider.Collection
	BBox          [4]float64
	Srid          int64
}

func NewPostgis(config provider.Config) (Postgis, error) {

	postgis := Postgis{}

	postgis.ApplicationId = config.ApplicationId
	postgis.UserVersion = config.UserVersion

	db, err := sqlx.Open("postgres", config.Datasource.PostGIS.Connection)

	postgis.Collections = config.Datasource.Collections
	postgis.BBox = config.Datasource.BBox
	postgis.Srid = int64(config.Datasource.Srid)

	if err != nil {
		return postgis, err
	}

	db.SetConnMaxLifetime(time.Minute * 15)

	postgis.db = db

	ctx := context.Background()

	postgis.UserVersion, _ = postgis.GetVersion(ctx, db)

	return postgis, nil
}

func (postgis Postgis) Close() error {
	return postgis.db.Close()
}

func (postgis Postgis) GetFeatures(ctx context.Context, db *sqlx.DB, collection provider.Collection, whereMap map[string]string, offset uint64, limit uint64, featureId interface{}, bbox [4]float64) (result FeatureCollectionGeoJSON, err error) {
	result = FeatureCollectionGeoJSON{}
	if len(bbox) > 4 {
		err = errors.New("bbox with 6 elements not supported")
		return
	}

	var FeatureIDColumn string

	if collection.Columns.Fid == "" {
		FeatureIDColumn = collection.Properties[0]
	} else {
		FeatureIDColumn = collection.Columns.Fid
	}

	tableName := fmt.Sprintf(`%s.%s`, collection.Schemaname, collection.Tablename)
	selectClause := fmt.Sprintf(`l."%s", st_asgeojson(st_forcesfs(l."%s")) as %s, l."%s"`, FeatureIDColumn, collection.Columns.Geometry, collection.Columns.Geometry, collection.Columns.Offset)

	// SELECT FEATURES
	for _, tf := range collection.Properties {
		if tf == collection.Columns.Geometry || tf == collection.Columns.Fid {
			continue
		}
		selectClause += fmt.Sprintf(`, l."%v"`, tf)
	}

	args := []interface{}{bbox[0], bbox[1], bbox[2], bbox[3], collection.Srid, offset, limit}

	additionalWhere := ""
	additionalWhereIndex := 8

	if collection.Filter != "" {
		additionalWhere += fmt.Sprintf(` %s AND `, collection.Filter)
	}

	if featureId != nil {
		additionalWhere = fmt.Sprintf(` l."%s"=$%d AND `, FeatureIDColumn, additionalWhereIndex)
		args = append(args, featureId)
		additionalWhereIndex++
	}
	// no JSONB features where clause as usual
	if !collection.Jsonb && len(whereMap) > 0 {
		for k := range whereMap {
			additionalWhere = fmt.Sprintf(` l."%s"=$%d AND `, k, additionalWhereIndex)
			args = append(args, featureId)
			additionalWhereIndex++
		}
	}
	// JSONB
	if collection.Jsonb && len(whereMap) > 0 {
		// JSONB COLUMN
		JSONBColumn := collection.Properties[0]
		//l."properties"@>  '{"lokaalID": "G1978.7afeb17a5c384f6bb08c2350e3f15b07"}'
		data, marshalErr := json.Marshal(whereMap)
		if marshalErr != nil {
			log.Printf("Could not marshal map %v", whereMap)
			err = marshalErr
			return
		}
		additionalWhere = fmt.Sprintf(` l."%s"@> '%s' AND `, JSONBColumn, string(data))

	}

	// query information with selection
	query := fmt.Sprintf(`SELECT %s FROM %s l WHERE %s st_intersects(st_makeenvelope($1,$2,$3,$4,$5), l."%s") AND l."%s" > $6 ORDER BY l."%s" LIMIT $7;`,
		selectClause, tableName, additionalWhere, collection.Columns.Geometry, collection.Columns.Offset, collection.Columns.Offset)
	rows, err := db.Queryx(query, args...)

	if err != nil {
		log.Printf("err during query: %v - %v", query, err)
		return
	}
	defer rowsClose(query, rows)

	cols, err := rows.Columns()
	if err != nil {
		return
	}

	result.NumberReturned = 0
	result.Type = "FeatureCollection"
	result.Features = make([]*Feature, 0)

	for rows.Next() {
		if err = ctx.Err(); err != nil {
			return
		}

		result.NumberReturned++

		vals := make([]interface{}, len(cols))
		valPtrs := make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++ {
			valPtrs[i] = &vals[i]
		}

		if err = rows.Scan(valPtrs...); err != nil {
			log.Printf("err reading row values: %v", err)
			return
		}

		feature := &Feature{Type: "Feature", Properties: make(map[string]interface{})}

		for i, colName := range cols {
			// check if the context cancelled or timed out
			if err = ctx.Err(); err != nil {
				return
			}

			//columnType := colTypes[i]
			if vals[i] == nil {
				continue
			}

			switch colName {
			case FeatureIDColumn:
				ID, err := provider.ConvertFeatureID(vals[i])
				if err != nil {
					return result, err
				}
				switch identifier := ID.(type) {
				case uint64:
					feature.ID = identifier
				case string:
					feature.ID = identifier
				}

			case collection.Columns.Offset:
				ofsset, ok := vals[i].(int64)
				if !ok {
					//log.Printf("unexpected column type for geom field. got %t", vals[i])
					return result, errors.New("unexpected column type for offset field. expected int")
				}
				result.Offset = ofsset

			case collection.Columns.Geometry:

				geomData, ok := vals[i].(string)
				if !ok {
					//log.Printf("unexpected column type for geom field. got %t", vals[i])
					return result, errors.New("unexpected column type for geom field. expected blob")
				}

				geometry := geojson.Geometry{}

				err := json.Unmarshal([]byte(geomData), &geometry)
				if err != nil {
					return result, err
				}
				feature.Geometry = geometry

			case collection.Columns.BBox:
				// Skip these columns used for bounding box and zoom filtering
				continue
			case "properties": // predefined jsonb
				switch v := vals[i].(type) {
				case []uint8:
					asBytes := make([]byte, len(v))
					for j := 0; j < len(v); j++ {
						asBytes[j] = v[j]
					}
					feature.Properties = make(map[string]interface{})
					err := json.Unmarshal(asBytes, &feature.Properties)
					if err != nil {
						return result, err
					}

				}
			default:
				// Grab any non-nil, non-id, non-bounding box, & non-geometry column as a tag
				switch v := vals[i].(type) {
				case []uint8:
					asBytes := make([]byte, len(v))
					for j := 0; j < len(v); j++ {
						asBytes[j] = v[j]
					}
					feature.Properties[colName] = string(asBytes)
				case int64:
					feature.Properties[colName] = v
				case float64:
					feature.Properties[colName] = v
				case time.Time:
					feature.Properties[colName] = v
				case string:
					feature.Properties[colName] = v
				case bool:
					feature.Properties[colName] = v
				default:
					log.Printf("unexpected type for postgis column data: %v: %T", cols[i], v)
				}
			}
		}
		result.Features = append(result.Features, feature)
	}

	return
}

func (postgis *Postgis) GetVersion(ctx context.Context, db *sqlx.DB) (string, error) {

	if postgis.UserVersion != "" {
		return postgis.UserVersion, nil
	}

	query := "SELECT PostGIS_full_version() as v;"
	// retrieve
	rows, err := db.Queryx(query)

	if err != nil {
		log.Printf("err during query: %v - %v", query, err)
		return "", err
	}
	defer rowsClose(query, rows)

	for rows.Next() {
		if err = ctx.Err(); err != nil {
			return "", err
		}
		row := struct {
			V string `db:"v"`
		}{}
		err := rows.StructScan(&row)

		if err != nil {
			log.Printf("err during query: %v - %v", query, err)
			return "", err
		}

		postgis.UserVersion = row.V
	}

	return postgis.UserVersion, nil
}

func rowsClose(query string, rows *sqlx.Rows) {

	err := rows.Close()

	if err != nil {
		log.Printf("err during closing rows: %v - %v", query, err)
	}

}
