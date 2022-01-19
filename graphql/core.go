package graphql

import (
	"net/http"
	"oaf-server/codegen"
	"oaf-server/core"
)

// NewGetLandingPageProvider passes the request to the Core NewGetLandingPageProvider with the GeoPackage Config
func (gp *GraphqlProvider) NewGetLandingPageProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetLandingPageProvider(gp.Config.Service)(r)
}

// NewGetApiProvider passes the request to the Core NewGetApiProvider with the GeoPackage Config
func (gp *GraphqlProvider) NewGetApiProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetApiProvider(gp.ApiProcessed)(r)
}

// NewGetConformanceDeclarationProvider passes the request to the Core NewGetConformanceDeclarationProvider with the GeoPackage Config
func (gp *GraphqlProvider) NewGetConformanceDeclarationProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetConformanceDeclarationProvider(gp.ApiProcessed)(r)
}

// NewGetCollectionsProvider passes the request to the Core NewGetCollectionsProvider with the GeoPackage Config
func (gp *GraphqlProvider) NewGetCollectionsProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetCollectionsProvider(gp.Config)(r)
}

// NewDescribeCollectionProvider passes the request to the Core NewDescribeCollectionProvider with the GeoPackage Config
func (gp *GraphqlProvider) NewDescribeCollectionProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewDescribeCollectionProvider(gp.Config)(r)
}
