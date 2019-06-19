package provider_postgis

import (
	"net/http"
	. "wfs3_server/codegen"
	pc "wfs3_server/provider_common"
)

type GetLandingPageProvider struct {
	Links []Link `json:"links"`
}

func (provider *PostgisProvider) NewGetLandingPageProvider(r *http.Request) (Provider, error) {
	return pc.NewGetLandingPageProvider(provider.serviceEndpoint)(r)
}
