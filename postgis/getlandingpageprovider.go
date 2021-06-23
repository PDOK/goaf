package postgis

import (
	"net/http"
	"oaf-server/codegen"
	"oaf-server/provider"
)

func (pp *PostgisProvider) NewGetLandingPageProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewGetLandingPageProvider(pp.CommonProvider.ServiceEndpoint)(r)
}
