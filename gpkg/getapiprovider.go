package gpkg

import (
	"log"
	"net/http"
	"oaf-server/codegen"
	"oaf-server/provider"

	"github.com/getkin/kin-openapi/openapi3"
)

type GetApiProvider struct {
	data        *openapi3.T
	contenttype string
}

func (gp *GeoPackageProvider) NewGetApiProvider(r *http.Request) (codegen.Provider, error) {
	p := &GetApiProvider{}

	ct, err := provider.GetContentType(r, p.String())
	if err != nil {
		return nil, err
	}
	p.contenttype = ct

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

func (gap *GetApiProvider) ContentType() string {
	return gap.contenttype
}

func (gap *GetApiProvider) String() string {
	return "api"
}

func (gap *GetApiProvider) SrsId() string {
	return "n.a"
}
