package provider_gpkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	. "wfs3_server/codegen"
)

type GetLandingPageProvider struct {
	data []Link
}

func (provider *GeoPackageProvider) NewGetLandingPageProvider(r *http.Request) (Provider, error) {

	ct := r.Header.Get("Content-Type")

	p := &GetLandingPageProvider{}

	if ct == "" {
		ct = JSONContentType
	}

	links, _ := provider.createLinks(fmt.Sprintf("%s/", provider.ServerEndpoint), "self", ct)
	apiLink, _ := provider.createLinks(fmt.Sprintf("%s/api", provider.ServerEndpoint), "service", ct)                     // /api, "service", ct)
	conformanceLink, _ := provider.createLinks(fmt.Sprintf("%s/conformance", provider.ServerEndpoint), "conformance", ct) // /conformance, "conformance", ct)
	dataLink, _ := provider.createLinks(fmt.Sprintf("%s/collections", provider.ServerEndpoint), "data", ct)               // /collections, "collections", ct)

	p.data = append(p.data, links...)
	p.data = append(p.data, apiLink...)
	p.data = append(p.data, conformanceLink...)
	p.data = append(p.data, dataLink...)

	return p, nil
}

func (provider *GetLandingPageProvider) Provide() (interface{}, error) {
	return provider.data, nil
}

func (provider *GetLandingPageProvider) MarshalJSON(interface{}) ([]byte, error) {
	return json.Marshal(provider.data)
}
func (provider *GetLandingPageProvider) MarshalHTML(interface{}) ([]byte, error) {
	// todo create html template pdok
	return json.Marshal(provider.data)
}
