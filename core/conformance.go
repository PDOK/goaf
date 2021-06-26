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

type GetConformanceDeclarationProvider struct {
	data        map[string][]string
	contenttype string
}

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

func (gcdp *GetConformanceDeclarationProvider) Provide() (interface{}, error) {
	return gcdp.data, nil
}

func (gcdp *GetConformanceDeclarationProvider) ContentType() string {
	return gcdp.contenttype
}

func (gcdp *GetConformanceDeclarationProvider) String() string {
	return "getconformancedeclaration"
}

func (gcdp *GetConformanceDeclarationProvider) SrsId() string {
	return "n.a."
}
