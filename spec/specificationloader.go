package spec

import (
	"github.com/getkin/kin-openapi/openapi3"
	"io/ioutil"
	"log"
)

func GetSwagger(serviceSpecPath string) (*openapi3.Swagger, error) {

	yaml, err := ioutil.ReadFile(serviceSpecPath)
	if err != nil {
		log.Fatalf("Cannot find file %s", serviceSpecPath)
	}
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromYAMLData(yaml)
	// tweak for missing schemaś reference to geojson
	swagger.Components.Schemas["geometryGeoJSON"] = &openapi3.SchemaRef{Ref: "http://geojson.org/schema/Geometry.json"}
	swagger.Components.Schemas["featureGeoJSON"] = &openapi3.SchemaRef{Ref: "http://geojson.org/schema/Feature.json"}
	return swagger, err
}
