package core

import (
	"net/http"
	"oaf-server/codegen"

	"github.com/getkin/kin-openapi/openapi3"
)

// GetApiProvider is returned by the NewGetApiProvider
// containing the data and contenttype for the response
type GetApiProvider struct {
	data        *openapi3.T
	contenttype string
}

// NewGetApiProvider handles the request and return the GetApiProvider
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

// Provide provides the data
func (gap *GetApiProvider) Provide() (interface{}, error) {
	return gap.data, nil
}

// ContentType returns the ContentType
func (gap *GetApiProvider) ContentType() string {
	return gap.contenttype
}

// String returns the provider name
func (gap *GetApiProvider) String() string {
	return "api"
}

// SrsId returns the srsid
func (gap *GetApiProvider) SrsId() string {
	return "n.a"
}
