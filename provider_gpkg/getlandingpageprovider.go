package provider_gpkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	. "wfs3_server/codegen"
)

type GetLandingPageProvider struct {
	links []Link
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

	p.links = append(p.links, links...)
	p.links = append(p.links, apiLink...)
	p.links = append(p.links, conformanceLink...)
	p.links = append(p.links, dataLink...)

	return p, nil
}

func (provider *GetLandingPageProvider) Provide() (interface{}, error) {
	return provider, nil
}

func (provider *GetLandingPageProvider) MarshalJSON(interface{}) ([]byte, error) {
	return json.Marshal(provider)
}
func (provider *GetLandingPageProvider) MarshalHTML(interface{}) ([]byte, error) {
	// todo create html template pdok
	return json.Marshal(provider)
}
