package core

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNewGetLandingPageProvider(t *testing.T) {

	service := Service{}

	provider, err := NewGetLandingPageProvider(service)(&http.Request{URL: &url.URL{RawQuery: ""}, Header: map[string][]string{}})
	if err != nil {
		t.Errorf("NewGetLandingPageProvider(serviceEndpoint) = %v, want %v", err, nil)
	}

	provided, err := provider.Provide()
	if err != nil {
		t.Errorf("NewGetLandingPageProvider.Provide() = %v, want %v", err, nil)
	}
	_, ok := provided.(*GetLandingPageProvider)
	if !ok {
		t.Errorf("NewGetLandingPageProvider.Provide() has incorrect type = %v, want %v", reflect.ValueOf(provided).Type(), "*GetLandingPageProvider")
	}

	links := provided.(*GetLandingPageProvider).Links

	if len(links) != 11 {
		t.Errorf("GetLandingPageProvider.Provide() has incorrect number of links = %v, want %v", len(links), 11)
	}
}
