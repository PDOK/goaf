package provider_dummy

import (
	"net/http"
	"wfs3_server/codegen"
)

// Starter example with all interface methods in place
type DummyProviders struct {
}

func (provider *DummyProviders) Init() error {

	return nil
}

func (provider *DummyProviders) NewGetApiProvider(r *http.Request) (codegen.Provider, error) {

	return nil, nil
}

func (provider *DummyProviders) NewGetLandingPageProvider(r *http.Request) (codegen.Provider, error) {

	return nil, nil
}

func (provider *DummyProviders) NewDescribeCollectionsProvider(r *http.Request) (codegen.Provider, error) {

	return nil, nil
}

func (provider *DummyProviders) NewDescribeCollectionProvider(r *http.Request) (codegen.Provider, error) {
	return nil, nil
}

func (provider *DummyProviders) NewGetFeaturesProvider(r *http.Request) (codegen.Provider, error) {
	return nil, nil
}

func (provider *DummyProviders) NewGetFeatureProvider(r *http.Request) (codegen.Provider, error) {
	return nil, nil
}

func (provider *DummyProviders) NewGetRequirementsClassesProvider(r *http.Request) (codegen.Provider, error) {
	return nil, nil
}
