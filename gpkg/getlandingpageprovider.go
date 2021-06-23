package gpkg

import (
	"net/http"
	"oaf-server/codegen"
	pc "oaf-server/provider"
)

func (gp *GeoPackageProvider) NewGetLandingPageProvider(r *http.Request) (codegen.Provider, error) {
	return pc.NewGetLandingPageProvider(gp.CommonProvider.ServiceEndpoint)(r)
}
