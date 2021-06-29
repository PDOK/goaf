package postgis

import (
	"net/http"
	"oaf-server/codegen"
	"oaf-server/core"
)

// NewGetLandingPageProvider passes the request to the Core NewGetLandingPageProvider with the Postgis Config
func (pp *PostgisProvider) NewGetLandingPageProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetLandingPageProvider(pp.Config.Service)(r)
}

// NewGetApiProvider passes the request to the Core NewGetApiProvider with the Postgis Config
func (pp *PostgisProvider) NewGetApiProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetApiProvider(pp.ApiProcessed)(r)
}

// NewGetConformanceDeclarationProvider passes the request to the Core NewGetConformanceDeclarationProvider with the Postgis Config
func (pp *PostgisProvider) NewGetConformanceDeclarationProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetConformanceDeclarationProvider(pp.ApiProcessed)(r)
}

// NewGetCollectionsProvider passes the request to the Core NewGetCollectionsProvider with the Postgis Config
func (pp *PostgisProvider) NewGetCollectionsProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetCollectionsProvider(pp.Config)(r)
}

// NewDescribeCollectionProvider passes the request to the Core NewDescribeCollectionProvider with the Postgis Config
func (pp *PostgisProvider) NewDescribeCollectionProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewDescribeCollectionProvider(pp.Config)(r)
}
