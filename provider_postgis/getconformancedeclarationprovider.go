package provider_postgis

import (
	"encoding/json"
	"net/http"
	. "wfs3_server/codegen"
)

type GetConformanceDeclarationProvider struct {
	data []string
}

func (provider *PostgisProvider) NewGetConformanceDeclarationProvider(r *http.Request) (Provider, error) {

	ct := r.Header.Get("Content-Type")

	p := &GetConformanceDeclarationProvider{}

	if ct == "" {
		ct = JSONContentType
	}

	p.data = []string{"http://www.opengis.net/spec/wfs-1/3.0/req/core", "http://www.opengis.net/spec/wfs-1/3.0/req/geojson"}

	return p, nil
}

func (provider *GetConformanceDeclarationProvider) Provide() (interface{}, error) {
	return provider.data, nil
}

func (provider *GetConformanceDeclarationProvider) MarshalJSON(interface{}) ([]byte, error) {
	return json.Marshal(provider.data)
}
func (provider *GetConformanceDeclarationProvider) MarshalHTML(interface{}) ([]byte, error) {
	// todo create html template pdok
	return json.Marshal(provider.data)
}
