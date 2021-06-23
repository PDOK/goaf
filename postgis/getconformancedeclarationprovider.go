package postgis

import (
	"errors"
	"net/http"
	"oaf-server/codegen"
)

type GetConformanceDeclarationProvider struct {
	data []string
}

func (pp *PostgisProvider) NewGetConformanceDeclarationProvider(r *http.Request) (codegen.Provider, error) {
	p := &GetConformanceDeclarationProvider{}
	path := r.URL.Path
	pathItem := pp.ApiProcessed.Paths.Find(path)
	if pathItem == nil {
		return p, errors.New("Invalid path :" + path)
	}

	for k := range r.URL.Query() {
		if notfound := pathItem.Get.Parameters.GetByInAndName("query", k) == nil; notfound {
			return p, errors.New("Invalid query parameter :" + k)
		}
	}
	p.data = []string{"http://www.opengis.net/spec/wfs-1/3.0/req/core", "http://www.opengis.net/spec/wfs-1/3.0/req/geojson"}
	return p, nil
}

func (gcdp *GetConformanceDeclarationProvider) Provide() (interface{}, error) {
	return gcdp.data, nil
}

func (gcdp *GetConformanceDeclarationProvider) String() string {
	return "getconformancedeclaration"
}

func (gcdp *GetConformanceDeclarationProvider) SrsId() string {
	return "n.a."
}
