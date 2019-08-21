package provider_common

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewGetLandingPageProvider(t *testing.T) {

	provider, err := NewGetLandingPageProvider("serviceEndpoint")(&http.Request{})
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

	if len(links) != 4 {
		t.Errorf("GetLandingPageProvider.Provide() has incorrect number of links = %v, want %v", len(links), 4)
	}
	_, err = provider.MarshalJSON(provided)
	if err != nil {
		t.Errorf("NewGetApiProvider.MarshalJSON() has error %v, want %v", err, nil)
	}

	_, err = provider.MarshalHTML(provided)
	if err != nil {
		t.Errorf("NewGetApiProvider.MarshalHTML() has error %v, want %v", err, nil)
	}

}
