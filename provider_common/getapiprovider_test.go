package provider_common

import (
	"github.com/getkin/kin-openapi/openapi3"
	"net/http"
	"reflect"
	"testing"
)

func TestNewGetApiProviderWithIncorrectSpecReference(t *testing.T) {

	_, err := NewGetApiProvider("../unknown/wfs3.0.json")(&http.Request{})
	if err == nil {
		t.Errorf("NewGetApiProvider(../unknown/wfs3.0.json) = %v, want %v", err, nil)
	}
}

func TestNewGetApiProvider(t *testing.T) {

	provider, err := NewGetApiProvider("../spec/wfs3.0.json")(&http.Request{})
	if err != nil {
		t.Errorf("NewGetApiProvider(../spec/wfs3.0.json) = %v, want %v", err, nil)
	}

	provided, err := provider.Provide()
	if err != nil {
		t.Errorf("NewGetApiProvider.Provide() = %v, want %v", err, nil)
	}
	_, ok := provided.(*openapi3.Swagger);
	if !ok {
		t.Errorf("NewGetApiProvider.Provide() has incorrect type = %v, want %v", reflect.ValueOf(provided).Type(), "*openapi3.Swagger")
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
