package provider_gpkg

import (
	"encoding/json"
	"github.com/go-spatial/geom/encoding/geojson"
	"net/http"
	cg "wfs3_server/codegen"
)

type GetFeatureProvider struct {
	data geojson.Feature
}

func (provider *GeoPackageProvider) NewGetFeatureProvider(r *http.Request) (cg.Provider, error) {

	collectionId, featureId := cg.ParametersForGetFeature(r)

	featureIdParam := provider.parseUint(featureId, 0)
	bboxParam := provider.GeoPackage.DefaultBBox

	ct := r.Header.Get("Content-Type")

	p := &GetFeatureProvider{}

	if ct == "" {
		ct = cg.JSONContentType
	}

	for _, cn := range provider.GeoPackage.Layers {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		fcGeoJSON, err := provider.GeoPackage.GetFeatures(r.Context(), provider.GeoPackage.DB, cn, collectionId, 0, 1, featureIdParam, bboxParam)

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
