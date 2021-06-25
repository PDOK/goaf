package provider

import (
	"io/ioutil"
	"log"
	"oaf-server/codegen"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ApplicationId string `yaml:"applicationid,omitempty"`
	UserVersion   string `yaml:"userversion,omitempty"`

	Openapi             string `yaml:"openapi"`
	DefaultFeatureLimit int    `yaml:"defaultfeaturelimit"`
	MaxFeatureLimit     int    `yaml:"maxfeaturelimit"`
	Datasource          Datasource
	Service             Service `yaml:"service" json:"service"`
}

type ContactPoint struct {
	Type        string `yaml:"@type" json:"@type,omitempty"`
	Email       string `yaml:"email" json:"email,omitempty"`
	Telephone   string `yaml:"telephone" json:"telephone,omitempty"`
	Url         string `yaml:"url" json:"url,omitempty"`
	ContactType string `yaml:"contactType" json:"contactType,omitempty"`
	Description string `yaml:"description" json:"description,omitempty"`
}
type Address struct {
	Type            string `yaml:"@type" json:"@type,omitempty"`
	StreetAddress   string `yaml:"streetAddress" json:"streetAddress,omitempty"`
	PostalCode      string `yaml:"postalCode" json:"postalCode,omitempty"`
	AddressLocality string `yaml:"addressLocality" json:"addressLocality,omitempty"`
	AddressRegion   string `yaml:"addressRegion" json:"addressRegion,omitempty"`
	AddressCountry  string `yaml:"addressCountry" json:"addressCountry,omitempty"`
}

type Provider struct {
	Type         string        `yaml:"@type" json:"@type"`
	Name         string        `yaml:"name" json:"name"`
	Url          string        `yaml:"url" json:"url"`
	Address      *Address      `yaml:"address" json:"address,omitempty"`           // pointer, omitting when empty
	ContactPoint *ContactPoint `yaml:"contactPoint" json:"contactPoint,omitempty"` // pointer, omitting when empty
}

type Service struct {
	Context     string   `yaml:"@context" json:"@context"`
	Type        string   `yaml:"@type" json:"@type"`
	Id          string   `yaml:"@id" json:"@id"`
	Url         string   `yaml:"url" json:"url"`
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	Keywords    []string `yaml:"keywords" json:"keywords"`
	License     string   `yaml:"license" json:"license"`
	LicenseName string   `yaml:"licenseName"` // do not output field to json
	Provider    Provider `yaml:"provider" json:"provider"`
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

	Links []codegen.Link `yaml:"links"`
}

type Columns struct {
	Fid      string `yaml:"fid"`
	Offset   string `yaml:"offset"`
	BBox     string `yaml:"bbox"`
	Geometry string `yaml:"geometry"`
}

func NewService() Service {
	address := Address{
		Type: "PostalAddress",
	}
	contactPoint := ContactPoint{
		Type: "Contactpoint",
	}
	provider := Provider{
		Type:         "Organization",
		ContactPoint: &contactPoint,
		Address:      &address,
	}
	service := Service{
		Context:  "https://schema.org/",
		Type:     "DataCatalog",
		Provider: provider,
	}
	return service
}

func (c *Config) ReadConfig(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read file from path (%v) with error: %v", path, err)
	}

	c.Service = NewService()

	yaml.Unmarshal(bytes, c)

	c.Service.Id = c.Service.Url
	add := c.Service.Provider.Address

	if add.AddressCountry == "" && add.PostalCode == "" && add.AddressLocality == "" && add.AddressRegion == "" && add.StreetAddress == "" {
		c.Service.Provider.Address = nil
	}

	cp := c.Service.Provider.ContactPoint

	if cp.ContactType == "" && cp.Description == "" && cp.Email == "" && cp.Telephone == "" && cp.Url == "" {
		c.Service.Provider.ContactPoint = nil
	}

	c.Service.Provider.Address = nil

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

	if c.Service.Url == "" {
		c.Service.Url = "http://localhost:8080"
	}

}
