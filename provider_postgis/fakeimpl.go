package provider_postgis

import (
	"net/http"
	"wfs3_server/codegen"
	_ "wfs3_server/codegen"
)

func (provider *PostgisProvider) NewGetApiProvider(r *http.Request) (codegen.Provider, error) {

	return nil, nil
}

func (provider *PostgisProvider) NewGetLandingPageProvider(r *http.Request) (codegen.Provider, error) {

	return nil, nil
}

func (provider *PostgisProvider) NewDescribeCollectionsProvider(r *http.Request) (codegen.Provider, error) {

	return nil, nil
}

func (provider *PostgisProvider) NewDescribeCollectionProvider(r *http.Request) (codegen.Provider, error) {

	return nil, nil
}

func (provider *PostgisProvider) NewGetFeaturesProvider(r *http.Request) (codegen.Provider, error) {

	return nil, nil
}

func (provider *PostgisProvider) NewGetFeatureProvider(r *http.Request) (codegen.Provider, error) {

	return nil, nil
}

func (provider *PostgisProvider) NewGetRequirementsClassesProvider(r *http.Request) (codegen.Provider, error) {

	return nil, nil
}
