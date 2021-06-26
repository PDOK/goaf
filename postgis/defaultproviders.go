package postgis

import (
	"net/http"
	"oaf-server/codegen"
	"oaf-server/provider"
)

func (pp *PostgisProvider) NewDescribeCollectionProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewDescribeCollectionProvider(pp.Config)(r)
}

func (pp *PostgisProvider) NewGetApiProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewGetApiProvider(pp.ApiProcessed)(r)
}

func (pp *PostgisProvider) NewGetConformanceDeclarationProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewGetConformanceDeclarationProvider(pp.ApiProcessed)(r)
}

func (pp *PostgisProvider) NewGetLandingPageProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewGetLandingPageProvider(pp.Config.Service)(r)
}

func (pp *PostgisProvider) NewGetCollectionsProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewGetCollectionsProvider(pp.Config)(r)
}
