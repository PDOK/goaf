package provider_postgis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-spatial/geom/encoding/geojson"
	"github.com/jmoiron/sqlx"
	"log"
	"regexp"
	"strconv"
	"time"
)

// mandatory according to geopackage specification
const (
	metatable_gpkg_contents        = "gpkg_contents"
	metatable_gpkg_spatial_ref_sys = " gpkg_spatial_ref_sys"
)

type PostgisLayer struct {
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

type Postgis struct {
	ApplicationId string
	UserVersion   string
	DB            *sqlx.DB
	Layers        []PostgisLayer
	DefaultBBox   []float64
	SrsId         int64
}

func NewPostgis(connectionStr string, featureTables []string) (Postgis, error) {

	postgis := &Postgis{}

	// Get all feature tables
	db, err := sqlx.Open("postgres", connectionStr)
	if err != nil {
		return *postgis, err
	}

	postgis.DB = db

	ctx := context.Background()

	applicationId, _ := postgis.GetApplicationID(ctx, db)
	version, _ := postgis.GetVersion(ctx, db)

	layers, err := postgis.GetLayers(ctx, db)

	log.Printf("| GEOPACKAGE DETAILS \n")
	log.Printf("|\n")
	log.Printf("| 	FILE: %s, APPLICATION: %s, VERSION: %d", connectionStr, applicationId, version)
	log.Printf("|\n")
	log.Printf("| 	NUMBER OF LAYERS: %d", len(layers))
	log.Printf("|\n")
	// determine query bbox
	for i, layer := range layers {
		log.Printf("| 	LAYER: %d. ID: %s, SRS_ID: %d, TABLE: %s PK: %s, FEATURES : %v\n", i+1, layer.Identifier, layer.SrsId, layer.Features[0], layer.Features[1], layer.Features[2:])

		if i == 0 {
			postgis.DefaultBBox = []float64{layer.MinX, layer.MinY, layer.MaxX, layer.MaxY}
			postgis.SrsId = layer.SrsId
		}
		if layer.MinX < postgis.DefaultBBox[0] {
			postgis.DefaultBBox[0] = layer.MinX
		}
		if layer.MinY < postgis.DefaultBBox[1] {
			postgis.DefaultBBox[1] = layer.MinY
		}
		if layer.MaxX > postgis.DefaultBBox[2] {
			postgis.DefaultBBox[2] = layer.MaxX
		}
		if layer.MaxY > postgis.DefaultBBox[3] {
			postgis.DefaultBBox[3] = layer.MaxY
		}
	}
	log.Printf("| \n")
	log.Printf("| 	BBOX: [%f,%f,%f,%f], SRS_ID:%d", postgis.DefaultBBox[0], postgis.DefaultBBox[1], postgis.DefaultBBox[2], postgis.DefaultBBox[3], postgis.SrsId)

	return *postgis, nil
}

func (gpkg *Postgis) Close() error {
	return gpkg.DB.Close()
}

func (gpkg *Postgis) GetLayers(ctx context.Context, db *sqlx.DB) (result []PostgisLayer, err error) {

	if gpkg.Layers != nil {
		result = gpkg.Layers
		err = nil
		return
	}

	re := regexp.MustCompile(`\"(.*?)\"|'(.*?)'`)

	query := `SELECT
			  c.table_name, c.data_type, c.identifier, c.description, c.last_change, c.min_x, c.min_y, c.max_x, c.max_y, c.srs_id, gc.column_name, gc.geometry_type_name, sm.sql
			  FROM
			  gpkg_contents c JOIN gpkg_geometry_columns gc ON c.table_name == gc.table_name JOIN sqlite_master sm ON c.table_name = sm.tbl_name
		      WHERE
			  c.data_type = 'features' AND sm.type = 'table' AND c.min_x IS NOT NULL`

	rows, err := db.Queryx(query)
	defer rows.Close()

	if err != nil {
		log.Printf("err during query: %v - %v", query, err)
		return
	}

	gpkg.Layers = make([]PostgisLayer, 0)

	for rows.Next() {
		if err = ctx.Err(); err != nil {
			return
		}
		row := PostgisLayer{}
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

func (gpkg Postgis) GetFeatures(ctx context.Context, db *sqlx.DB, layer PostgisLayer, collectionId string, offset uint64, limit uint64, featureId uint64, bbox []float64) (result FeatureCollectionGeoJSON, err error) {
	// Features bit of a hack // layer.Features => tablename, PK, ...FEATURES, assuming create table in sql statement first is PK
	result = FeatureCollectionGeoJSON{}
	if len(bbox) > 4 {
		err = errors.New("bbox with 6 elements not supported!")
		return
	}

	rtreeTablenName := fmt.Sprintf("rtree_%s_%s", layer.TableName, layer.ColumnName)
	selectClause := fmt.Sprintf("l.`%s`, l.`%s`", layer.Features[1], layer.ColumnName)

	for _, tf := range layer.Features[2:] { // [2:] skip tablename and PK
		if tf == layer.ColumnName {
			continue
		}
		selectClause += fmt.Sprintf(", l.`%v`", tf)
	}

	additionalWhere := ""
	if featureId > 0 { // explicit count should be 1
		additionalWhere = fmt.Sprintf(" `id`=%d AND ", featureId)
	} else {

		// count total with selection
		queryCount := fmt.Sprintf("SELECT count(*) AS `total` FROM `%s` WHERE minx <= %v AND maxx >= %v AND miny <= %v AND maxy >= %v;",
			rtreeTablenName, bbox[2], bbox[0], bbox[3], bbox[1])
		count, err := db.Query(queryCount)
		if err != nil {
			log.Printf("err during query: %v - %v", queryCount, err)
			return result, err
		}
		defer count.Close()

		if count.Next() {
			err = count.Scan(&result.NumberMatched)
			if err != nil {
				log.Printf("err reading row values: %v", err)
				return result, err
			}
		}
	}

	// query information with selection
	query := fmt.Sprintf("SELECT %s FROM `%s` l WHERE l.`%s` IN (SELECT id FROM %s WHERE %s minx <= %v AND maxx >= %v AND miny <= %v AND maxy >= %v ORDER BY ID LIMIT %d OFFSET %d);",
		selectClause, layer.TableName, layer.Features[1], rtreeTablenName, additionalWhere, bbox[2], bbox[0], bbox[3], bbox[1], limit, offset)
	rows, err := db.Queryx(query)

	if err != nil {
		log.Printf("err during query: %v - %v", query, err)
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return
	}

	result.NumberReturned = 0
	result.Type = "FeatureCollection"
	result.Features = make([]geojson.Feature, 0)

	for rows.Next() {
		if err = ctx.Err(); err != nil {
			return
		}

		if featureId > 0 {
			additionalWhere = fmt.Sprintf(" l.`%s`=%d AND ", layer.Features[1], featureId)
			result.NumberMatched++
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

		feature := geojson.Feature{Properties: make(map[string]interface{})}

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
			case layer.Features[1]:
				ID, err := convertFeatureID(vals[i])
				if err != nil {
					return result, err
				}
				feature.ID = &ID
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
				default:
					log.Printf("unexpected type for sqlite column data: %v: %T", cols[i], v)
				}
			}
		}
		result.Features = append(result.Features, feature)
	}

	return
}

func (gpkg *Postgis) GetApplicationID(ctx context.Context, db *sqlx.DB) (string, error) {

	if gpkg.ApplicationId != "" {
		return gpkg.ApplicationId, nil
	}

	gpkg.ApplicationId = "POSTGIS"

	return gpkg.ApplicationId, nil

}

func (postgis *Postgis) GetVersion(ctx context.Context, db *sqlx.DB) (string, error) {

	if postgis.UserVersion != "" {
		return postgis.UserVersion, nil
	}

	query := "SELECT PostGIS_full_version();"
	// retrieve
	_, rows, err := executeRaw(ctx, db, query)
	if err != nil {
		log.Printf("err during query: %v - %v", query, err)
		return "", err
	}

	if len(rows) == 0 {
		return "", errors.New("cannot determine postgis version")
	}

	postgis.UserVersion = rows[0][0].(string)

	return postgis.UserVersion, nil
}

func executeRaw(ctx context.Context, db *sqlx.DB, query string) (cols []string, rows [][]interface{}, err error) {

	rowz, err := db.Query(query)
	defer rowz.Close()

	if err != nil {
		log.Printf("err during query: %v - %v", query, err)
		return
	}

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

// convertFeatureID attempts to convert an interface value to an uint64
// copied from https://github.com/go-spatial/jivan
func convertFeatureID(v interface{}) (uint64, error) {
	switch aval := v.(type) {
	case float64:
		return uint64(aval), nil
	case int64:
		return uint64(aval), nil
	case uint64:
		return aval, nil
	case uint:
		return uint64(aval), nil
	case int8:
		return uint64(aval), nil
	case uint8:
		return uint64(aval), nil
	case uint16:
		return uint64(aval), nil
	case int32:
		return uint64(aval), nil
	case uint32:
		return uint64(aval), nil
	case string:
		return strconv.ParseUint(aval, 10, 64)
	default:
		return 0, errors.New(fmt.Sprintf("Cannot convert to numeric ID : %v", aval))
	}
}
