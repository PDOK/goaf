package provider_postgis

import (
	"net/http"
	cg "oaf-server/codegen"
	pc "oaf-server/provider_common"
)

func (provider *PostgisProvider) NewGetLandingPageProvider(r *http.Request) (cg.Provider, error) {
	return pc.NewGetLandingPageProvider(provider.CommonProvider.ServiceEndpoint)(r)
}
