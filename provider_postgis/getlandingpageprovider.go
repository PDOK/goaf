package provider_postgis

import (
	"encoding/json"
	"fmt"
	"net/http"
	. "wfs3_server/codegen"
)

type GetLandingPageProvider struct {
	Links []Link `json:"links"`
}

func (provider *PostgisProvider) NewGetLandingPageProvider(r *http.Request) (Provider, error) {

	ct := r.Header.Get("Content-Type")

	p := &GetLandingPageProvider{}

	if ct == "" {
		ct = JSONContentType
	}

	links, _ := provider.createLinks(fmt.Sprintf("%s/", provider.serviceEndpoint), "self", ct)
	apiLink, _ := provider.createLinks(fmt.Sprintf("%s/api", provider.serviceEndpoint), "service", ct)                     // /api, "service", ct)
	conformanceLink, _ := provider.createLinks(fmt.Sprintf("%s/conformance", provider.serviceEndpoint), "conformance", ct) // /conformance, "conformance", ct)
	dataLink, _ := provider.createLinks(fmt.Sprintf("%s/collections", provider.serviceEndpoint), "data", ct)               // /collections, "collections", ct)

	p.Links = append(p.Links, links...)
	p.Links = append(p.Links, apiLink...)
	p.Links = append(p.Links, conformanceLink...)
	p.Links = append(p.Links, dataLink...)

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
