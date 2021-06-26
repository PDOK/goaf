package provider

import (
	"net/http"
	"oaf-server/codegen"

	"github.com/getkin/kin-openapi/openapi3"
)

type GetApiProvider struct {
	data        *openapi3.T
	contenttype string
}

func NewGetApiProvider(api *openapi3.T) func(r *http.Request) (codegen.Provider, error) {

	return func(r *http.Request) (codegen.Provider, error) {
		p := &GetApiProvider{}

		ct, err := GetContentType(r, p.String())
		if err != nil {
			return nil, err
		}

		p.contenttype = ct
		p.data = api

		return p, nil
	}
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
