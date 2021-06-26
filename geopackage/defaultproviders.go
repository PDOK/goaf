package geopackage

import (
	"net/http"
	"oaf-server/codegen"
	"oaf-server/core"
)

func (gp *GeoPackageProvider) NewDescribeCollectionProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewDescribeCollectionProvider(gp.Config)(r)
}

func (gp *GeoPackageProvider) NewGetApiProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetApiProvider(gp.ApiProcessed)(r)
}

func (gp *GeoPackageProvider) NewGetConformanceDeclarationProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetConformanceDeclarationProvider(gp.ApiProcessed)(r)
}

func (gp *GeoPackageProvider) NewGetLandingPageProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetLandingPageProvider(gp.Config.Service)(r)
}

func (gp *GeoPackageProvider) NewGetCollectionsProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetCollectionsProvider(gp.Config)(r)
}
