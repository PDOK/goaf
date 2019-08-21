package provider_postgis

import (
	"errors"
	"fmt"
	"net/http"
	cg "wfs3_server/codegen"
)

type GetFeatureProvider struct {
	data Feature
}

func (provider *PostgisProvider) NewGetFeatureProvider(r *http.Request) (cg.Provider, error) {

	collectionId, featureId := cg.ParametersForGetFeature(r)

	featureIdParam := featureId
	bboxParam := provider.PostGis.BBox

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

		fcGeoJSON, err := provider.PostGis.GetFeatures(r.Context(), provider.PostGis.db, cn, collectionId, 0, 1, featureIdParam, bboxParam)

		if err != nil {
			return nil, err
		}

		if len(fcGeoJSON.Features) == 1 {
			p.data = fcGeoJSON.Features[0]
		} else {
			return p, errors.New(fmt.Sprintf("Feature with id: %s not found", string(featureId)))
		}

		break
	}

	return p, nil
}

func (provider *GetFeatureProvider) Provide() (interface{}, error) {
	return provider.data, nil
}

func (provider *GetFeatureProvider) String() string {
	return "getfeature"
}