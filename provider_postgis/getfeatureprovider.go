package provider_postgis

import (
	"encoding/json"
	"net/http"
	cg "wfs3_server/codegen"
)

type GetFeatureProvider struct {
	data Feature
}

func (provider *PostgisProvider) NewGetFeatureProvider(r *http.Request) (cg.Provider, error) {

	collectionId, featureId := cg.ParametersForGetFeature(r)

	featureIdParam := featureId
	bboxParam := provider.PostGis.DefaultBBox

	ct := r.Header.Get("Content-Type")

	p := &GetFeatureProvider{}

	if ct == "" {
		ct = cg.JSONContentType
	}

	for _, cn := range provider.PostGis.Layers {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		fcGeoJSON, err := provider.PostGis.GetFeatures(r.Context(), provider.PostGis.DB, cn, collectionId, 0, 1, featureIdParam, bboxParam)

		if err != nil {
			return nil, err
		}

		if len(fcGeoJSON.Features) == 1 {
			p.data = fcGeoJSON.Features[0]
		}

		break
	}

	return p, nil
}

func (provider *GetFeatureProvider) Provide() (interface{}, error) {
	return provider.data, nil
}

func (provider *GetFeatureProvider) MarshalJSON(interface{}) ([]byte, error) {
	return json.Marshal(provider.data)
}
func (provider *GetFeatureProvider) MarshalHTML(interface{}) ([]byte, error) {
	// todo create html template pdok
	return json.Marshal(provider.data)
}
