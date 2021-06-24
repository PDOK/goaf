package gpkg

import (
	"net/http"
	"oaf-server/codegen"
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

func (gp *GeoPackageProvider) NewGetConformanceDeclarationProvider(r *http.Request) (codegen.Provider, error) {

	p := &GetConformanceDeclarationProvider{}
	p.contenttype = r.Header.Get("Content-Type")

	d := make(map[string][]string)
	d["conformsTo"] = []string{core, oas30, html, gjson}

	p.data = d

	return p, nil
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
	return "n.a"
}
