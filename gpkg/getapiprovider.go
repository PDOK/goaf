package gpkg

import (
	"net/http"
	"oaf-server/codegen"
	"oaf-server/provider"
)

func (gp *GeoPackageProvider) NewGetApiProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewGetApiProvider(gp.ApiProcessed)(r)
}
