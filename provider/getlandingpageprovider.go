package provider

import (
	"fmt"
	"net/http"
	"oaf-server/codegen"
)

type GetLandingPageProvider struct {
	Links       []codegen.Link `json:"links"`
	contenttype string
}

func NewGetLandingPageProvider(serviceEndpoint string) func(r *http.Request) (codegen.Provider, error) {

	return func(r *http.Request) (codegen.Provider, error) {

		ct := r.Header.Get("Content-Type")

		p := &GetLandingPageProvider{}
		p.contenttype = ct

		links, _ := CreateLinks("landing page", serviceEndpoint, "self", ct)
		apiLink, _ := GetApiLinks(fmt.Sprintf("%s/api", serviceEndpoint))                                                    // /api, "service", ct)
		conformanceLink, _ := CreateLinks("capabilities", fmt.Sprintf("%s/conformance", serviceEndpoint), "conformance", ct) // /conformance, "conformance", ct)
		dataLink, _ := CreateLinks("collections", fmt.Sprintf("%s/collections", serviceEndpoint), "data", ct)                // /collections, "collections", ct)

		p.Links = append(p.Links, links...)
		p.Links = append(p.Links, apiLink...)
		p.Links = append(p.Links, conformanceLink...)
		p.Links = append(p.Links, dataLink...)

		return p, nil
	}
}

func (glp *GetLandingPageProvider) Provide() (interface{}, error) {
	return glp, nil
}

func (glp *GetLandingPageProvider) ContentType() string {
	return glp.contenttype
}

func (glp *GetLandingPageProvider) String() string {
	return "landingpage"
}

func (glp *GetLandingPageProvider) SrsId() string {
	return "n.a"
}
