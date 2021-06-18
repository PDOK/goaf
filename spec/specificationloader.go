package spec

import (
	"log"

	"github.com/getkin/kin-openapi/openapi3"
)

func GetOpenAPI(serviceSpecPath string) (*openapi3.T, error) {

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	openapi, err := loader.LoadFromFile(serviceSpecPath)
	if err != nil {
		log.Printf("Cannot Loadswagger from file :%s", serviceSpecPath)
	}
	// tweak for missing schemaś reference to geojson
	//swagger.Components.Schemas["geometryGeoJSON"] = &openapi3.SchemaRef{Ref: "http://geojson.org/schema/Geometry.json"}
	//swagger.Components.Schemas["featureGeoJSON"] = &openapi3.SchemaRef{Ref: "http://geojson.org/schema/Feature.json"}
	return openapi, err
}
