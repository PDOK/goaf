package provider_gpkg

import (
	"log"
	"net/http"
	"wfs3_server/codegen"

	"github.com/getkin/kin-openapi/openapi3"
)

type GetApiProvider struct {
	data *openapi3.T
}

func (provider *GeoPackageProvider) NewGetApiProvider(r *http.Request) (codegen.Provider, error) {
	p := &GetApiProvider{}
	var err error
	if provider.Api == nil {
		log.Printf("Could not get Swagger Specification")
		return p, err
	}

	p.data = provider.Api
	return p, nil
}

func (provider *GetApiProvider) Provide() (interface{}, error) {
	return provider.data, nil
}

func (provider *GetApiProvider) String() string {
	return "api"
}

func (provider *GetApiProvider) SrsId() string {
	return "n.a"
}
