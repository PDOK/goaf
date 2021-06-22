package provider_common

import (
	"fmt"
	"net/http"
	cg "oaf-server/codegen"
)

type GetLandingPageProvider struct {
	Links []cg.Link `json:"links"`
}

func NewGetLandingPageProvider(serviceEndpoint string) func(r *http.Request) (cg.Provider, error) {

	return func(r *http.Request) (cg.Provider, error) {

		ct := r.Header.Get("Content-Type")

		p := &GetLandingPageProvider{}

		links, _ := CreateLinks("landing page", serviceEndpoint, "self", ct)
		apiLink, _ := CreateLinks("openapi3 specification", fmt.Sprintf("%s/api", serviceEndpoint), "service", ct)           // /api, "service", ct)
		conformanceLink, _ := CreateLinks("capabilities", fmt.Sprintf("%s/conformance", serviceEndpoint), "conformance", ct) // /conformance, "conformance", ct)
		dataLink, _ := CreateLinks("collections", fmt.Sprintf("%s/collections", serviceEndpoint), "data", ct)                // /collections, "collections", ct)

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

func (provider *GetLandingPageProvider) String() string {
	return "landingpage"
}

func (provider *GetLandingPageProvider) SrsId() string {
	return "n.a"
}
