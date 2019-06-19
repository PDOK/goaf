package provider_common

import (
	"encoding/json"
	"fmt"
	"net/http"
	. "wfs3_server/codegen"
)

type GetLandingPageProvider struct {
	Links []Link `json:"links"`
}

func NewGetLandingPageProvider(serviceEndpoint string) func(r *http.Request) (Provider, error) {

	return func(r *http.Request) (Provider, error) {

		ct := r.Header.Get("Content-Type")

		p := &GetLandingPageProvider{}

		if ct == "" {
			ct = JSONContentType
		}

		links, _ := CreateLinks(fmt.Sprintf("%s/", serviceEndpoint), "self", ct)
		apiLink, _ := CreateLinks(fmt.Sprintf("%s/api", serviceEndpoint), "service", ct)                     // /api, "service", ct)
		conformanceLink, _ := CreateLinks(fmt.Sprintf("%s/conformance", serviceEndpoint), "conformance", ct) // /conformance, "conformance", ct)
		dataLink, _ := CreateLinks(fmt.Sprintf("%s/collections", serviceEndpoint), "data", ct)               // /collections, "collections", ct)

		p.Links = append(p.Links, links...)
		p.Links = append(p.Links, apiLink...)
		p.Links = append(p.Links, conformanceLink...)
		p.Links = append(p.Links, dataLink...)

		return p, nil
	}
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
