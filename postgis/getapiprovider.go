package postgis

import (
	"net/http"
	"oaf-server/codegen"
	"oaf-server/provider"
)

func (pp *PostgisProvider) NewGetApiProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewGetApiProvider(pp.ApiProcessed)(r)
}
