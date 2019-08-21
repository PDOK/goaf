package provider_common

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3"
	"log"
	"net/http"
	"wfs3_server/codegen"
	"wfs3_server/spec"
)

type GetApiProvider struct {
	data *openapi3.Swagger
}

var swagger *openapi3.Swagger

func NewGetApiProvider(serviceSpecPath string) func(r *http.Request) (codegen.Provider, error) {

	return func(r *http.Request) (codegen.Provider, error) {

		ct := r.Header.Get("Content-Type")
		p := &GetApiProvider{}

		if ct == "" {
			ct = codegen.JSONContentType
		}

		var err error
		if swagger == nil {
			swagger, err = spec.GetSwagger(serviceSpecPath)
			if err != nil {
				log.Printf("Error parsing swagger file %s", serviceSpecPath)
				return p, err
			}
		}

		p.data = swagger

		return p, nil
	}
}

func (provider *GetApiProvider) Provide() (interface{}, error) {
	return provider.data, nil
}

func (provider *GetApiProvider) MarshalJSON(interface{}) ([]byte, error) {
	return json.Marshal(provider.data)
}
func (provider *GetApiProvider) MarshalHTML(interface{}) ([]byte, error) {
	// todo html swagger API templating etc. with custom item viewer?
	return json.Marshal(provider.data)
}
