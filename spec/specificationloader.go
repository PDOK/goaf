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

	return openapi, err
}
