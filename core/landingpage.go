package core

import (
	"fmt"
	"net/http"
	"oaf-server/codegen"
)

// GetLandingPageProvider is returned by the NewGetLandingPageProvider
// containing the data and contenttype for the response
type GetLandingPageProvider struct {
	Links       []codegen.Link `json:"links,omitempty"`
	contenttype string
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	*Service
}

// NewGetLandingPageProvider handles the request and return the GetLandingPageProvider
func NewGetLandingPageProvider(serviceConfig Service) func(r *http.Request) (codegen.Provider, error) {

	return func(r *http.Request) (codegen.Provider, error) {
		p := &GetLandingPageProvider{}
		reqContentType, err := GetContentType(r, p.String())

		if err != nil {
			return nil, err
		}

		p.contenttype = reqContentType
		p.Service = &serviceConfig

		links, _ := CreateLinks("landing page", p.String(), serviceConfig.Url, "self", reqContentType)
		apiLink, _ := GetApiLinks(fmt.Sprintf("%s/api", serviceConfig.Url))                                                                            // /api, "service", ct)
		conformanceLink, _ := CreateLinks("capabilities", p.String(), fmt.Sprintf("%s/conformance", serviceConfig.Url), "conformance", reqContentType) // /conformance, "conformance", ct)
		dataLink, _ := CreateLinks("collections", p.String(), fmt.Sprintf("%s/collections", serviceConfig.Url), "data", reqContentType)                // /collections, "collections", ct)

		p.Links = append(p.Links, links...)
		p.Links = append(p.Links, apiLink...)
		p.Links = append(p.Links, conformanceLink...)
		p.Links = append(p.Links, dataLink...)

		if p.contenttype == "application/json" {
			p.Title = p.Service.Name
			p.Description = p.Service.Description
			p.Service = nil
		} else if p.contenttype == "application/ld+json" {
			p.Links = nil
		}
		return p, nil
	}
}

// Provide provides the data
func (glp *GetLandingPageProvider) Provide() (interface{}, error) {
	return glp, nil
}

// ContentType returns the ContentType
func (glp *GetLandingPageProvider) ContentType() string {
	return glp.contenttype
}

// String returns the provider name
func (glp *GetLandingPageProvider) String() string {
	return "landingpage"
}

// SrsId returns the srsid
func (glp *GetLandingPageProvider) SrsId() string {
	return "n.a"
}
