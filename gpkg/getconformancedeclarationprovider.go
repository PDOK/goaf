package gpkg

import (
	"net/http"
	"oaf-server/codegen"
	"oaf-server/provider"
)

type GetConformanceDeclarationProvider struct {
	data                  []string
	contenttype           string
	supportedContentTypes map[string]string
}

func (gp *GeoPackageProvider) NewGetConformanceDeclarationProvider(r *http.Request) (codegen.Provider, error) {

	p := &GetConformanceDeclarationProvider{}

	ct, err := provider.GetContentType(r, p.ProviderType())
	if err != nil {
		return nil, err
	}

	p.contenttype = ct

	p.data = []string{"http://www.opengis.net/spec/wfs-1/3.0/req/core", "http://www.opengis.net/spec/wfs-1/3.0/req/geojson"}

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

func (gcdp *GetConformanceDeclarationProvider) ProviderType() string {
	return provider.CapabilitesProvider
}
