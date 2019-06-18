package provider_postgis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-spatial/geom/encoding/geojson"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

// mandatory according to geopackage specification
const (
	metatable_gpkg_contents        = "gpkg_contents"
	metatable_gpkg_spatial_ref_sys = " gpkg_spatial_ref_sys"
)

type IdNotFoundError struct {
	err string
}

type PostgisLayer struct {
	SchemaName      string `yaml:"SchemaName,omitempty"`
	TableName       string `yaml:"TableName,omitempty"`
	Description     string `yaml:"Description,omitempty"`
	Identifier      string `yaml:"Identifier,omitempty"`
	GeometryColumn  string `yaml:"GeometryColumn,omitempty"`
	GeometryType    string `yaml:"GeometryType,omitempty"`
	FeatureIDColumn string `yaml:"FeatureIDColumn,omitempty"`
	OffsetColumn    string `yaml:"OffsetColumn,omitempty"`

	BBoxGeometryColumn string `yaml:"BBoxGeometryColumn,omitempty"`

	BBox     []float64 `yaml:"BBox,omitempty"`
	SrsId    int64     `yaml:"SrsId,omitempty"`
	Features []string  `yaml:"Features,omitempty"`
}

type Postgis struct {
	ConnectionStr string `yaml:"ConnectionStr,omitempty"`
	ApplicationId string `yaml:"ApplicationId,omitempty"`
	UserVersion   string `yaml:"UserVersion,omitempty"`
	db            *sqlx.DB
	Layers        []PostgisLayer `yaml:"Layers,omitempty"`
	BBox          []float64      `yaml:"BBox,omitempty"`
	SrsId         int64          `yaml:"SrsId,omitempty"`
}

func NewPostgis(configfilePath string) (Postgis, error) {

	postgis := Postgis{}

	configFile, err := ioutil.ReadFile(configfilePath)

	if err != nil {
		log.Println("Could not read config file: %s, using default CRS Map", configFile)
		return postgis, err
	} else {
		err := yaml.Unmarshal(configFile, &postgis)

		if err != nil {
			log.Println("Could not unmarshal config file: %s", configfilePath)
			return postgis, err
		}

	}

	db, err := sqlx.Open("postgres", postgis.ConnectionStr)
	if err != nil {
		return postgis, err
	}

	postgis.db = db

	ctx := context.Background()

	postgis.UserVersion, _ = postgis.GetVersion(ctx, db)

	return postgis, nil
}

func (gpkg *Postgis) Close() error {
	return gpkg.db.Close()
}

func (postgis Postgis) GetFeatures(ctx context.Context, db *sqlx.DB, layer PostgisLayer, collectionId string, offset uint64, limit uint64, featureId interface{}, bbox []float64) (result FeatureCollectionGeoJSON, err error) {
	result = FeatureCollectionGeoJSON{}
	if len(bbox) > 4 {
		err = errors.New("bbox with 6 elements not supported!")
		return
	}

	var FeatureIDColumn string

	if layer.FeatureIDColumn == "" {
		FeatureIDColumn = layer.Features[0]
	} else {
		FeatureIDColumn = layer.FeatureIDColumn
	}

	tablenName := fmt.Sprintf(`%s.%s`, layer.SchemaName, layer.TableName)
	selectClause := fmt.Sprintf(`l."%s", st_asgeojson(l."%s") as %s`, FeatureIDColumn, layer.GeometryColumn, layer.GeometryColumn)

	for _, tf := range layer.Features {

		if tf == layer.GeometryColumn || tf == FeatureIDColumn {
			continue
		}
		selectClause += fmt.Sprintf(`, l."%v"`, tf)
	}

	additionalWhere := ""

	if featureId != nil {
		switch identifier := featureId.(type) {
		case uint64:
			additionalWhere = fmt.Sprintf(` l."%s"=%d AND `, FeatureIDColumn, identifier)
		case string:
			additionalWhere = fmt.Sprintf(` l."%s"='%s' AND `, FeatureIDColumn, identifier)
		}
	}

	// query information with selection
	query := fmt.Sprintf(`SELECT %s FROM %s l WHERE %s st_intersects(st_makeenvelope(%v,%v,%v,%v, %v), l."%s") AND l."%s" > %d ORDER BY l."%s" LIMIT %d;`,
		selectClause, tablenName, additionalWhere, bbox[0], bbox[1], bbox[2], bbox[3], layer.SrsId, layer.GeometryColumn, layer.OffsetColumn, offset, layer.OffsetColumn, limit)

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
			case FeatureIDColumn:
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

			case layer.GeometryColumn:

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

			case layer.OffsetColumn, layer.BBoxGeometryColumn:
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
	case string:
		return aval, nil

	default:
		return 0, errors.New(fmt.Sprintf("Cannot convert ID : %v", aval))
	}
}
