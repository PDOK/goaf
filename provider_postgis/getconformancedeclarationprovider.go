package provider_postgis

import (
	"errors"
	"net/http"
	cg "wfs3_server/codegen"
)

type GetConformanceDeclarationProvider struct {
	data []string
}

func (provider *PostgisProvider) NewGetConformanceDeclarationProvider(r *http.Request) (cg.Provider, error) {
	p := &GetConformanceDeclarationProvider{}
	path := r.URL.Path
	pathItem := provider.ApiProcessed.Paths.Find(path)
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

func (provider *GetConformanceDeclarationProvider) Provide() (interface{}, error) {
	return provider.data, nil
}

func (provider *GetConformanceDeclarationProvider) String() string {
	return "getconformancedeclaration"
}

func (provider *GetConformanceDeclarationProvider) SrsId() string {
	return "n.a."
}
