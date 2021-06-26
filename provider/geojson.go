package provider

import (
	"oaf-server/codegen"

	"github.com/go-spatial/geom/encoding/geojson"
)

type FeatureCollectionGeoJSON struct {
	NumberReturned int64          `json:"numberReturned,omitempty"`
	TimeStamp      string         `json:"timeStamp,omitempty"`
	Type           string         `json:"type"`
	Features       []*Feature     `json:"features"`
	Links          []codegen.Link `json:"links,omitempty"`
	NumberMatched  int64          `json:"numberMatched,omitempty"`
	Crs            string         `json:"crs,omitempty"`
	Offset         int64          `json:"-"`
}

type Feature struct {
	// overwrite ID in geojson.Feature so strings are also allowed as id
	ID interface{} `json:"id,omitempty"`
	geojson.Feature
	// Added Links in de document
	Links []codegen.Link `json:"links,omitempty"`
}
