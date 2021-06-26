package core

import (
	"errors"
	"net/http"
	"oaf-server/codegen"

	"github.com/getkin/kin-openapi/openapi3"
)

const (
	core  = "http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/core"
	oas30 = "http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/oas30"
	html  = "http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/html"
	gjson = "http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/geojson"
)

// GetConformanceDeclarationProvider is returned by the NewGetConformanceDeclarationProvider
// containing the data and contenttype for the response
type GetConformanceDeclarationProvider struct {
	data        map[string][]string
	contenttype string
}

// NewGetConformanceDeclarationProvider handles the request and return the GetConformanceDeclarationProvider
func NewGetConformanceDeclarationProvider(api *openapi3.T) func(r *http.Request) (codegen.Provider, error) {
	return func(r *http.Request) (codegen.Provider, error) {
		p := &GetConformanceDeclarationProvider{}

		ct, err := GetContentType(r, p.String())

		if err != nil {
			return nil, err
		}

		p.contenttype = ct
		path := r.URL.Path
		pathItem := api.Paths.Find(path)
		if pathItem == nil {
			return p, errors.New("Invalid path :" + path)
		}

		for k := range r.URL.Query() {
			if notfound := pathItem.Get.Parameters.GetByInAndName("query", k) == nil; notfound {
				return p, errors.New("Invalid query parameter :" + k)
			}
		}

		d := make(map[string][]string)
		d["conformsTo"] = []string{core, oas30, html, gjson}

		p.data = d
		return p, nil
	}

}

// Provide provides the data
func (gcdp *GetConformanceDeclarationProvider) Provide() (interface{}, error) {
	return gcdp.data, nil
}

// ContentType returns the ContentType
func (gcdp *GetConformanceDeclarationProvider) ContentType() string {
	return gcdp.contenttype
}

// String returns the provider name
func (gcdp *GetConformanceDeclarationProvider) String() string {
	return "conformance"
}

// SrsId returns the srsid
func (gcdp *GetConformanceDeclarationProvider) SrsId() string {
	return "n.a."
}
