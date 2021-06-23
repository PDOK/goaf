package gpkg

import (
	"log"
	"net/http"
	"oaf-server/codegen"

	"github.com/getkin/kin-openapi/openapi3"
)

type GetApiProvider struct {
	data *openapi3.T
}

func (gp *GeoPackageProvider) NewGetApiProvider(r *http.Request) (codegen.Provider, error) {
	p := &GetApiProvider{}
	var err error
	if gp.Api == nil {
		log.Printf("Could not get Swagger Specification")
		return p, err
	}

	p.data = gp.Api
	return p, nil
}

func (gap *GetApiProvider) Provide() (interface{}, error) {
	return gap.data, nil
}

func (gap *GetApiProvider) String() string {
	return "api"
}

func (gap *GetApiProvider) SrsId() string {
	return "n.a"
}
