package postgis

import (
	"fmt"
	"net/http"
	"oaf-server/codegen"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type GetApiProvider struct {
	data *openapi3.T
}

func (pp *PostgisProvider) NewGetApiProvider(r *http.Request) (codegen.Provider, error) {
	p := &GetApiProvider{}

	p.data = pp.ApiProcessed
	return p, nil
}

func CreateProvidesSpecificParameters(provider *PostgisProvider) *openapi3.T {

	api := provider.Api

	copy := &openapi3.T{
		OpenAPI:      api.OpenAPI,
		Info:         api.Info,
		Servers:      api.Servers,
		Paths:        make(map[string]*openapi3.PathItem),
		Components:   api.Components,
		Security:     api.Security,
		ExternalDocs: api.ExternalDocs,
	}

	copy.Components.Extensions = nil

	delete(copy.Components.Parameters, "collectionId")

	for k, v := range provider.Api.Paths {
		if !strings.Contains(k, "{collectionId}") {
			v.Extensions = nil
			copy.Paths[k] = v
		}
	}

	// adjust swagger to accommodate individual parameters
	for _, layer := range provider.PostGis.Layers {
		for k, v := range provider.Api.Paths {
			if strings.Contains(k, "{collectionId}") {
				k := strings.Replace(k, "{collectionId}", strings.ToLower(layer.TableName), 1)
				params := openapi3.NewParameters()
				paramsQueryExists := false

				for _, p := range v.Get.Parameters {
					if strings.Contains(p.Ref, "collectionId") {
						continue
					}

					if p.Value.Name != "collectionId" {
						params = append(params, p)
						if p.Value.In == "query" {
							paramsQueryExists = true
						}
					}
				}
				// only add vendor specific parameters to query params are already allowed
				if paramsQueryExists {
					for _, specificParam := range layer.VendorSpecificParameters {
						sp := openapi3.NewQueryParameter(specificParam)
						sp.Description = fmt.Sprintf("Vendor specific parameter : %s", specificParam)
						sp.Required = false
						sp.Schema = &openapi3.SchemaRef{
							Ref: "",
							Value: &openapi3.Schema{
								Type: "object",
							},
						}
						params = append(params, &openapi3.ParameterRef{
							Ref:   "#/components/parameters/" + specificParam,
							Value: sp,
						})

						copy.Components.Parameters[specificParam] = &openapi3.ParameterRef{
							Value: sp,
						}
					}
				}

				copy.Paths[k] = v
				copy.Paths[k].Get.Parameters = params
				copy.Paths[k].Get.Extensions = nil

			}
		}
	}
	return copy
}

func (gap *GetApiProvider) Provide() (interface{}, error) {
	return gap.data, nil
}

func (gap *GetApiProvider) String() string {
	return "api"
}

func (gap *GetApiProvider) SrsId() string {
	return "n.a"
}
