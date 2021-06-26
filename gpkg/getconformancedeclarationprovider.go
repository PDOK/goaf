package gpkg

import (
	"net/http"
	"oaf-server/codegen"
	"oaf-server/provider"
)

func (gp *GeoPackageProvider) NewGetConformanceDeclarationProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewGetConformanceDeclarationProvider(gp.ApiProcessed)(r)
}
