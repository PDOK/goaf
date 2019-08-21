package spec

import (
	"github.com/getkin/kin-openapi/openapi3"
	"log"
)

func GetSwagger(serviceSpecPath string) (*openapi3.Swagger, error) {

	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true

	swagger, err := loader.LoadSwaggerFromFile(serviceSpecPath)
	if err != nil {
		log.Fatalf("Cannot Loadswagger from file :%s", serviceSpecPath)
	}
	// tweak for missing schemaś reference to geojson
	//swagger.Components.Schemas["geometryGeoJSON"] = &openapi3.SchemaRef{Ref: "http://geojson.org/schema/Geometry.json"}
	//swagger.Components.Schemas["featureGeoJSON"] = &openapi3.SchemaRef{Ref: "http://geojson.org/schema/Feature.json"}
	return swagger, err
}
