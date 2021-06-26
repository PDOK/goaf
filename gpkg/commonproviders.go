package gpkg

import (
	"net/http"
	"oaf-server/codegen"
	"oaf-server/provider"
)

func (gp *GeoPackageProvider) NewDescribeCollectionProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewDescribeCollectionProvider(gp.Config)(r)
}

func (gp *GeoPackageProvider) NewGetApiProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewGetApiProvider(gp.ApiProcessed)(r)
}

func (gp *GeoPackageProvider) NewGetConformanceDeclarationProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewGetConformanceDeclarationProvider(gp.ApiProcessed)(r)
}

func (gp *GeoPackageProvider) NewGetLandingPageProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewGetLandingPageProvider(gp.Config.Service)(r)
}
