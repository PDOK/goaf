package provider_postgis

import (
	"net/http"
	cg "wfs3_server/codegen"
	pc "wfs3_server/provider_common"
)

func (provider *PostgisProvider) NewGetLandingPageProvider(r *http.Request) (cg.Provider, error) {
	return pc.NewGetLandingPageProvider(provider.CommonProvider.ServiceEndpoint)(r)
}
