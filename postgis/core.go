package postgis

import (
	"net/http"
	"oaf-server/codegen"
	"oaf-server/core"
)

func (pp *PostgisProvider) NewDescribeCollectionProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewDescribeCollectionProvider(pp.Config)(r)
}

func (pp *PostgisProvider) NewGetApiProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetApiProvider(pp.ApiProcessed)(r)
}

func (pp *PostgisProvider) NewGetConformanceDeclarationProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetConformanceDeclarationProvider(pp.ApiProcessed)(r)
}

func (pp *PostgisProvider) NewGetLandingPageProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetLandingPageProvider(pp.Config.Service)(r)
}

func (pp *PostgisProvider) NewGetCollectionsProvider(r *http.Request) (codegen.Provider, error) {
	return core.NewGetCollectionsProvider(pp.Config)(r)
}
