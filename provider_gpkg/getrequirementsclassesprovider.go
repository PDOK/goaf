package provider_gpkg

import (
	"encoding/json"
	"net/http"
	. "wfs3_server/codegen"
)

type GetRequirementsClassesProvider struct {
	data []string
}

func (provider *GeoPackageProvider) NewGetRequirementsClassesProvider(r *http.Request) (Provider, error) {

	ct := r.Header.Get("Content-Type")

	p := &GetRequirementsClassesProvider{}

	if ct == "" {
		ct = JSONContentType
	}

	p.data = []string{"http://www.opengis.net/spec/wfs-1/3.0/req/core", "http://www.opengis.net/spec/wfs-1/3.0/req/geojson"}

	return p, nil
}

func (provider *GetRequirementsClassesProvider) Provide() (interface{}, error) {
	return provider.data, nil
}

func (provider *GetRequirementsClassesProvider) MarshalJSON(interface{}) ([]byte, error) {
	return json.Marshal(provider.data)
}
func (provider *GetRequirementsClassesProvider) MarshalHTML(interface{}) ([]byte, error) {
	// todo create html template pdok
	return json.Marshal(provider.data)
}
