package provider_gpkg

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	pc "oaf-server/provider_common"
	"os"
	"regexp"
	"time"

	"github.com/go-spatial/geom/encoding/geojson"
	"github.com/jmoiron/sqlx"
)

type GeoPackageLayer struct {
	TableName    string    `db:"table_name"`
	DataType     string    `db:"data_type"`
	Identifier   string    `db:"identifier"`
	Description  string    `db:"description"`
	ColumnName   string    `db:"column_name"`
	GeometryType string    `db:"geometry_type_name"`
	LastChange   time.Time `db:"last_change"`
	// bbox
	MinX     float64  `db:"min_x"`
	MinY     float64  `db:"min_y"`
	MaxX     float64  `db:"max_x"`
	MaxY     float64  `db:"max_y"`
	SrsId    int64    `db:"srs_id"`
	SQL      string   `db:"sql"`
	Features []string // first table, second PK, rest features
}

type GeoPackage struct {
	ApplicationId string
	UserVersion   int64
	DB            *sqlx.DB
	FeatureIdKey  string
	Layers        []GeoPackageLayer
	DefaultBBox   []float64
	SrsId         int64
}

func NewGeoPackage(filepath string, featureIdKey string) (GeoPackage, error) {

	gpkg := &GeoPackage{}

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return *gpkg, fmt.Errorf("geopackage invalid location : %s", filepath)
	}

	// Get all feature tables
	db, err := sqlx.Open("sqlite3", filepath)
	if err != nil {
		return *gpkg, err
	}

	gpkg.FeatureIdKey = featureIdKey
	gpkg.DB = db

	ctx := context.Background()

	applicationId, _ := gpkg.GetApplicationID(ctx, db)
	version, _ := gpkg.GetVersion(ctx, db)

	layers, _ := gpkg.GetLayers(ctx, db)

	log.Printf("| GEOPACKAGE DETAILS \n")
	log.Printf("|\n")
	log.Printf("| 	FILE: %s, APPLICATION: %s, VERSION: %d", filepath, applicationId, version)
	log.Printf("|\n")
	log.Printf("| 	NUMBER OF LAYERS: %d", len(layers))
	log.Printf("|\n")
	// determine query bbox
	for i, layer := range layers {
		log.Printf("| 	LAYER: %d. ID: %s, SRS_ID: %d, TABLE: %s PK: %s, FEATURES : %v\n", i+1, layer.Identifier, layer.SrsId, layer.Features[0], layer.Features[1], layer.Features[2:])

		if i == 0 {
			gpkg.DefaultBBox = []float64{layer.MinX, layer.MinY, layer.MaxX, layer.MaxY}
			gpkg.SrsId = layer.SrsId
		}
		if layer.MinX < gpkg.DefaultBBox[0] {
			gpkg.DefaultBBox[0] = layer.MinX
		}
		if layer.MinY < gpkg.DefaultBBox[1] {
			gpkg.DefaultBBox[1] = layer.MinY
		}
		if layer.MaxX > gpkg.DefaultBBox[2] {
			gpkg.DefaultBBox[2] = layer.MaxX
		}
		if layer.MaxY > gpkg.DefaultBBox[3] {
			gpkg.DefaultBBox[3] = layer.MaxY
		}
	}
	log.Printf("| \n")
	log.Printf("| 	BBOX: [%f,%f,%f,%f], SRS_ID:%d", gpkg.DefaultBBox[0], gpkg.DefaultBBox[1], gpkg.DefaultBBox[2], gpkg.DefaultBBox[3], gpkg.SrsId)

	return *gpkg, nil
}

func (gpkg *GeoPackage) Close() error {
	return gpkg.DB.Close()
}

func (gpkg *GeoPackage) GetLayers(ctx context.Context, db *sqlx.DB) (result []GeoPackageLayer, err error) {

	if gpkg.Layers != nil {
		result = gpkg.Layers
		err = nil
		return
	}

	re := regexp.MustCompile(`"(.*?)"|'(.*?)'`)

	query := `SELECT
			  c.table_name, c.data_type, c.identifier, c.description, c.last_change, c.min_x, c.min_y, c.max_x, c.max_y, c.srs_id, gc.column_name, gc.geometry_type_name, sm.sql
			  FROM
			  gpkg_contents c JOIN gpkg_geometry_columns gc ON c.table_name == gc.table_name JOIN sqlite_master sm ON c.table_name = sm.tbl_name
		      WHERE
			  c.data_type = 'features' AND sm.type = 'table' AND c.min_x IS NOT NULL`

	rows, err := db.Queryx(query)
	if err != nil {
		log.Printf("err during query: %v - %v", query, err)
		return
	}
	defer rowsClose(query, rows)

	gpkg.Layers = make([]GeoPackageLayer, 0)

	for rows.Next() {
		if err = ctx.Err(); err != nil {
			return
		}
		row := GeoPackageLayer{}
		err := rows.StructScan(&row)
		if err != nil {
			log.Fatalln(err)
		}

		row.Features = make([]string, 0)
		matches := re.FindAllStringSubmatch(row.SQL, -1)
		for _, match := range matches {
			row.Features = append(row.Features, match[1])
		}

		gpkg.Layers = append(gpkg.Layers, row)
	}

	result = gpkg.Layers

	return
}

func (gpkg GeoPackage) GetFeatures(ctx context.Context, db *sqlx.DB, layer GeoPackageLayer, collectionId string, offset uint64, limit uint64, featureId interface{}, bbox []float64) (result *FeatureCollectionGeoJSON, err error) {
	// Features bit of a hack // layer.Features => tablename, PK, ...FEATURES, assuming create table in sql statement first is PK
	result = &FeatureCollectionGeoJSON{}
	if len(bbox) > 4 {
		err = errors.New("bbox with 6 elements not supported")
		return
	}

	var featureIdKey string

	if gpkg.FeatureIdKey == "" {
		featureIdKey = layer.Features[1]
	} else {
		featureIdKey = gpkg.FeatureIdKey
	}

	rtreeTablenName := fmt.Sprintf("rtree_%s_%s", layer.TableName, layer.ColumnName)
	selectClause := fmt.Sprintf("l.`%s`, l.`%s`", featureIdKey, layer.ColumnName)

	for _, tf := range layer.Features[1:] { // [2:] skip tablename and PK
		if tf == layer.ColumnName || tf == featureIdKey {
			continue
		}
		selectClause += fmt.Sprintf(", l.`%v`", tf)
	}

	additionalWhere := ""

	if featureId != nil {
		additionalWhere = fmt.Sprintf(` l."%s"=$1 AND `, featureIdKey)
	}

	query := fmt.Sprintf("SELECT %s FROM `%s` l INNER JOIN `%s` g ON g.`id` = l.`fid` WHERE %s minx <= $2 AND maxx >= $3 AND miny <= $4 AND maxy >= $5 ORDER BY l.`%s` LIMIT $6 OFFSET $7;",
		selectClause, layer.TableName, rtreeTablenName, additionalWhere, featureIdKey)

	var rows *sqlx.Rows
	if featureId != nil {
		rows, err = db.Queryx(query, featureId, bbox[2], bbox[0], bbox[3], bbox[1], limit, offset)
	} else {
		rows, err = db.Queryx(query, bbox[2], bbox[0], bbox[3], bbox[1], limit, offset)
	}

	if err != nil {
		log.Printf("err during query: %v - %v", query, err)
		return
	}
	defer rowsClose(query, rows)
	cols, err := rows.Columns()

	if err != nil {
		log.Printf("err during query: %v - %v", query, err)
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
			case featureIdKey:
				ID, err := pc.ConvertFeatureID(vals[i])
				if err != nil {
					return result, err
				}
				switch identifier := ID.(type) {
				case uint64:
					feature.ID = identifier
				case string:
					feature.ID = identifier
				}

			case layer.ColumnName:

				geomData, ok := vals[i].([]byte)
				if !ok {
					//log.Printf("unexpected column type for geom field. got %t", vals[i])
					return result, errors.New("unexpected column type for geom field. expected blob")
				}

				_, geo, err := DecodeGeometry(geomData)
				if err != nil {
					return result, err
				}
				feature.Geometry = geojson.Geometry{Geometry: geo}

			case "minx", "miny", "maxx", "maxy", "min_zoom", "max_zoom":
				// Skip these columns used for bounding box and zoom filtering
				continue

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
					log.Printf("unexpected type for sqlite column data: %v: %T", cols[i], v)
				}
			}
		}
		result.Features = append(result.Features, feature)
	}

	return
}

func (gpkg *GeoPackage) GetApplicationID(ctx context.Context, db *sqlx.DB) (string, error) {

	if gpkg.ApplicationId != "" {
		return gpkg.ApplicationId, nil
	}

	query := "PRAGMA application_id"
	// retrieve
	_, rows, err := executeRaw(ctx, db, query)
	if err != nil {
		log.Printf("err during query: %v - %v", query, err)
		return "", err
	}

	if len(rows) == 0 {
		return "", errors.New("cannot determine geopackage application id")
	}

	// check length rows/colums
	applicationId := rows[0][0].(int64)

	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(applicationId))

	gpkg.ApplicationId = string(b[4:]) // should result in GPKG

	return gpkg.ApplicationId, nil

}

func (gpkg *GeoPackage) GetVersion(ctx context.Context, db *sqlx.DB) (int64, error) {

	if gpkg.UserVersion != 0 {
		return gpkg.UserVersion, nil
	}

	query := "PRAGMA user_version"
	// retrieve
	_, rows, err := executeRaw(ctx, db, query)
	if err != nil {
		log.Printf("err during query: %v - %v", query, err)
		return -1, err
	}
	// check length rows/colums
	if len(rows) == 0 {
		return 0, errors.New("cannot determine geopackage user_version")
	}

	gpkg.UserVersion = rows[0][0].(int64)

	return gpkg.UserVersion, nil
}

func executeRaw(ctx context.Context, db *sqlx.DB, query string) (cols []string, rows [][]interface{}, err error) {

	rowz, err := db.Query(query)

	if err != nil {
		log.Printf("err during query: %v - %v", query, err)
		return
	}
	defer func() {
		err := rowz.Close()
		if err != nil {
			log.Printf("err during closing rows: %v - %v", query, err)
		}
	}()

	cols, err = rowz.Columns()
	if err != nil {
		return
	}

	rows = make([][]interface{}, 0)

	for rowz.Next() {
		if err = ctx.Err(); err != nil {
			return
		}

		vals := make([]interface{}, len(cols))
		valPtrs := make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++ {
			valPtrs[i] = &vals[i]
		}

		if err = rowz.Scan(valPtrs...); err != nil {
			log.Printf("err reading row values: %v", err)
			return
		}

		row := make([]interface{}, len(cols))

		for i := range cols {
			// check if the context cancelled or timed out
			if err = ctx.Err(); err != nil {
				return
			}
			if vals[i] == nil {
				row[i] = nil
				continue
			}

			switch v := vals[i].(type) {
			case []uint8:
				asBytes := make([]byte, len(v))
				for j := 0; j < len(v); j++ {
					asBytes[j] = v[j]
				}
				row[i] = string(asBytes)
			case int64:
				//feature.Properties[cols[i]] = v
				row[i] = v
			default:
				log.Printf("unexpected type for sqlite column data: %v: %T", cols[i], v)
			}

			rows = append(rows, row)

		}

	}

	return
}

func rowsClose(query string, rows *sqlx.Rows) {

	err := rows.Close()

	if err != nil {
		log.Printf("err during closing rows: %v - %v", query, err)
	}

}
