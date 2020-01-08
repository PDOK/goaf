package provider_gpkg

import (
	"net/http"
	. "wfs3_server/codegen"
)

type GetConformanceDeclarationProvider struct {
	data []string
}

func (provider *GeoPackageProvider) NewGetConformanceDeclarationProvider(r *http.Request) (Provider, error) {

	p := &GetConformanceDeclarationProvider{}

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
	return "n.a"
}
