package gpkg

import (
	"net/http"
	"oaf-server/codegen"
)

type GetConformanceDeclarationProvider struct {
	data []string
}

func (gp *GeoPackageProvider) NewGetConformanceDeclarationProvider(r *http.Request) (codegen.Provider, error) {

	p := &GetConformanceDeclarationProvider{}

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
	return "n.a"
}
