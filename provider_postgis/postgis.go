package provider_postgis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-spatial/geom/encoding/geojson"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"strings"
	"time"
)

// mandatory according to geopackage specification
const (
	metatable_gpkg_contents        = "gpkg_contents"
	metatable_gpkg_spatial_ref_sys = " gpkg_spatial_ref_sys"
)

type PostgisLayer struct {
	SchemaName   string    `db:"table_schema"`
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
	FeatureIdKey  string
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

	layers, err := postgis.GetLayers(ctx, db, featureTables)

	log.Printf("| POSTGIS DETAILS \n")
	log.Printf("|\n")
	log.Printf("| 	CONNECTION: %s, APPLICATION: %s, VERSION: %d", connectionStr, applicationId, version)
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

func (postgis *Postgis) GetLayers(ctx context.Context, db *sqlx.DB, featureTables []string) (result []PostgisLayer, err error) {

	if postgis.Layers != nil {
		result = postgis.Layers
		err = nil
		return
	}

	queryTemplate := `SELECT "f_table_schema" as table_schema,
					 "f_table_name" as table_name, 
					 "f_geometry_column" as column_name,
					 "srid" as srs_id,
					 'features' as data_type,
					 "type" as geometry_type_name,
					 "f_table_name" || '_' || "f_geometry_column" as identifier
			  FROM public.geometry_columns where f_table_name in ('%s');`

	layerQuery := fmt.Sprintf(queryTemplate, strings.Join(featureTables, "','"))

	layers, err := db.Queryx(layerQuery)
	if err != nil {
		log.Printf("err during query: %v - %v", queryTemplate, err)
		return
	}
	defer layers.Close()

	postgis.Layers = make([]PostgisLayer, 0)

	for layers.Next() {
		if err = ctx.Err(); err != nil {
			return
		}
		layer := PostgisLayer{}
		err := layers.StructScan(&layer)
		if err != nil {
			log.Fatalln(err)
		}

		layer, err = extractQueryExtent(db, layer)
		if err != nil {
			log.Fatalln(err)
		}

		layer, err = extractFeature(db, layer)
		if err != nil {
			log.Fatalln(err)
		}

		postgis.Layers = append(postgis.Layers, layer)

	}

	result = postgis.Layers

	return
}

func extractFeature(db *sqlx.DB, layer PostgisLayer) (PostgisLayer, error) {
	type DbFeature struct {
		ColumnName string `db:"column_name"`
		DataType   string `db:"data_type"`
	}

	template := `SELECT c.column_name, c.data_type 
                 FROM information_schema.columns c
                 WHERE c.table_schema = '%s' AND c.table_name = '%s' 
                 ORDER BY c.ordinal_position;`

	query := fmt.Sprintf(template, layer.SchemaName, layer.TableName)
	q, err := db.Queryx(query)
	if err != nil {
		log.Printf("err during query : %v - %v", query, err)
		return layer, err
	}
	defer q.Close()

	layer.Features = make([]string, 0)
	layer.Features = append(layer.Features, layer.TableName)

	for q.Next() {

		feature := DbFeature{}
		err = q.StructScan(&feature)
		if err != nil {
			log.Printf("err during query : %v - %v", query, err)
			return layer, err
		}

		if feature.DataType == "USER-DEFINED" || feature.ColumnName == layer.ColumnName {
			continue
		}

		layer.Features = append(layer.Features, feature.ColumnName)
	}

	return layer, nil
}

func extractQueryExtent(db *sqlx.DB, layer PostgisLayer) (PostgisLayer, error) {
	template := `SELECT st_xmin(extent) as min_x, st_ymin(extent) as min_y, st_xmax(extent) as max_x, st_ymax(extent) as max_y FROM (SELECT ST_Extent(%s) AS extent FROM %s.%s) AS bbox;`
	query := fmt.Sprintf(template, layer.ColumnName, layer.SchemaName, layer.TableName)
	q, err := db.Queryx(query)
	if err != nil {
		log.Printf("err during query : %v - %v", query, err)
		return layer, err
	}
	defer q.Close()

	if q.Next() {
		err = q.StructScan(&layer)
		if err != nil {
			log.Printf("err during query : %v - %v", query, err)
			return layer, err
		}
	}

	return layer, nil
}

func (postgis Postgis) GetFeatures(ctx context.Context, db *sqlx.DB, layer PostgisLayer, collectionId string, offset uint64, limit uint64, featureId interface{}, bbox []float64) (result FeatureCollectionGeoJSON, err error) {
	result = FeatureCollectionGeoJSON{}
	if len(bbox) > 4 {
		err = errors.New("bbox with 6 elements not supported!")
		return
	}

	var featureIdKey string

	if postgis.FeatureIdKey == "" {
		featureIdKey = layer.Features[1]
	} else {
		featureIdKey = postgis.FeatureIdKey
	}

	tablenName := fmt.Sprintf(`%s.%s`, layer.SchemaName, layer.TableName)
	selectClause := fmt.Sprintf(`l."%s", st_asgeojson(st_forcesfs(l."%s")) as %s`, featureIdKey, layer.ColumnName, layer.ColumnName)

	for _, tf := range layer.Features[1:] {
		if tf == layer.ColumnName || tf == featureIdKey {
			continue
		}
		selectClause += fmt.Sprintf(`, l."%v"`, tf)
	}

	additionalWhere := ""

	if featureId != nil {
		switch identifier := featureId.(type) {
		case uint64:
			additionalWhere = fmt.Sprintf(` l."%s"=%d AND `, postgis.FeatureIdKey, identifier)
		case string:
			additionalWhere = fmt.Sprintf(` l."%s"="%s" AND `, postgis.FeatureIdKey, identifier)
		}
	}

	// query information with selection
	query := fmt.Sprintf(`SELECT %s FROM %s l WHERE %s st_intersects(st_makeenvelope(%v,%v,%v,%v, %v), st_forcesfs(l."%s")) ORDER BY l."%s" LIMIT %d OFFSET %d;`,
		selectClause, tablenName, additionalWhere, bbox[0], bbox[1], bbox[2], bbox[3], layer.SrsId, layer.ColumnName, featureIdKey, limit, offset)

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
	result.Features = make([]Feature, 0)

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

		feature := Feature{Type: "Feature", Properties: make(map[string]interface{})}

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
				ID, err := convertFeatureID(vals[i])
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
					log.Printf("unexpected type for postgis column data: %v: %T", cols[i], v)
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

	query := "SELECT PostGIS_full_version() as v;"
	// retrieve
	rows, err := db.Queryx(query)

	if err != nil {
		log.Printf("err during query: %v - %v", query, err)
		return "", err
	}
	defer rows.Close()

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
func convertFeatureID(v interface{}) (interface{}, error) {
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
	case []byte:
		return string(aval), nil
	default:
		return 0, errors.New(fmt.Sprintf("Cannot convert ID : %v", aval))
	}
}
