package postgis

import (
	"net/http"
	"oaf-server/codegen"
	"oaf-server/provider"
)

func (pp *PostgisProvider) NewGetConformanceDeclarationProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewGetConformanceDeclarationProvider(pp.ApiProcessed)(r)
}
