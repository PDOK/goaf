package provider

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ApplicationId string `yaml:"applicationid,omitempty"`
	UserVersion   string `yaml:"userversion,omitempty"`

	Endpoint            string `yaml:"endpoint"`
	Openapi             string `yaml:"openapi"`
	DefaultFeatureLimit int    `yaml:"defaultfeaturelimit"`
	MaxFeatureLimit     int    `yaml:"maxfeaturelimit"`
	Datasource          Datasource
}

type Datasource struct {
	Geopackage  *Geopackage  `yaml:"gpkg"`
	PostGIS     *PostGIS     `yaml:"postgis"`
	Collections []Collection `yaml:"collections"`
	BBox        [4]float64   `yaml:"bbox"`
	Srid        int          `yaml:"srid"`
}

type Geopackage struct {
	File string `yaml:"file"`
	Fid  string `yaml:"fid"`
}

type PostGIS struct {
	Connection string `yaml:"connection"`
}

type Collection struct {
	Schemaname  string `yaml:"schemaname"`
	Tablename   string `yaml:"tablename"`
	Identifier  string `yaml:"identifier"`
	Description string `yaml:"description"`
	Filter      string `yaml:"filter,omitempty"`

	Columns                  *Columns   `yaml:"columns"`
	Geometrytype             string     `yaml:"geometrytype,omitempty"`
	BBox                     [4]float64 `yaml:"bbox"`
	Srid                     int        `yaml:"srid"`
	VendorSpecificParameters []string   `yaml:"vendorspecificparameters"`
	Jsonb                    bool       `yaml:"jsonb"`
	Properties               []string   `yaml:"properties"`
}

type Columns struct {
	Fid      string `yaml:"fid"`
	Offset   string `yaml:"offset"`
	BBox     string `yaml:"bbox"`
	Geometry string `yaml:"geometry"`
}

func (c *Config) ReadConfig(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read file from path (%v) with error: %v", path, err)
	}

	yaml.Unmarshal(bytes, c)

	// set defaults if none are provided
	if c.DefaultFeatureLimit < 1 {
		c.DefaultFeatureLimit = 100
	}

	if c.MaxFeatureLimit < 1 {
		c.MaxFeatureLimit = 500
	}

	if c.Openapi == "" {
		c.Openapi = "spec/oaf.json"
	}

	if c.Endpoint == "" {
		c.Endpoint = "http://localhost:8080"
	}

}
