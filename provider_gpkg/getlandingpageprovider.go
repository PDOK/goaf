package provider_gpkg

import (
	"net/http"
	. "wfs3_server/codegen"
	"wfs3_server/provider_common"
)

type GetLandingPageProvider struct {
	Links []Link `json:"links"`
}

func (provider *GeoPackageProvider) NewGetLandingPageProvider(r *http.Request) (Provider, error) {
	return provider_common.NewGetLandingPageProvider(provider.serviceEndpoint)(r)
}
