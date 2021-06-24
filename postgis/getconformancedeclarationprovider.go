package postgis

import (
	"errors"
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

func (pp *PostgisProvider) NewGetConformanceDeclarationProvider(r *http.Request) (codegen.Provider, error) {
	p := &GetConformanceDeclarationProvider{}
	p.contenttype = r.Header.Get("Content-Type")
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
	return "n.a."
}
