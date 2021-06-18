package provider_postgis

import (
	"errors"
	"fmt"
	"net/http"
	cg "wfs3_server/codegen"
	pc "wfs3_server/provider_common"
)

type GetFeatureProvider struct {
	data  *Feature
	srsid string
}

func (provider *PostgisProvider) NewGetFeatureProvider(r *http.Request) (cg.Provider, error) {

	collectionId, featureId, _ := cg.ParametersForGetFeature(r)

	featureIdParam := featureId
	bboxParam := provider.PostGis.BBox

	p := &GetFeatureProvider{srsid: fmt.Sprintf("EPSG:%d", provider.PostGis.SrsId)}

	path := r.URL.Path
	ct := r.Header.Get("Content-Type")

	for _, cn := range provider.PostGis.Layers {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		pathItem := provider.ApiProcessed.Paths.Find("/collections/pand/items/{featureId}")
		if pathItem == nil {
			return p, errors.New("Invalid path :" + path)
		}

		for k := range r.URL.Query() {
			if notfound := pathItem.Get.Parameters.GetByInAndName("query", k) == nil; notfound {
				return p, errors.New("Invalid query parameter :" + k)
			}
		}

		whereMap := make(map[string]string)
		fcGeoJSON, err := provider.PostGis.GetFeatures(r.Context(), provider.PostGis.db, cn, whereMap, 0, 1, featureIdParam, bboxParam)

		if err != nil {
			return nil, err
		}

		if len(fcGeoJSON.Features) >= 1 {
			feature := fcGeoJSON.Features[0]

			hrefBase := fmt.Sprintf("%s%s", provider.CommonProvider.ServiceEndpoint, path) // /collections
			links, _ := pc.CreateLinks("feature", hrefBase, "self", ct)
			feature.Links = links

			p.data = feature

		} else {
			return p, fmt.Errorf("Feature with id: %s not found", string(featureId))
		}

		return p, nil
	}

	return p, errors.New("Cannot find layer : " + collectionId)
}

func (provider *GetFeatureProvider) Provide() (interface{}, error) {
	return provider.data, nil
}

func (provider *GetFeatureProvider) String() string {
	return "getfeature"
}

func (provider *GetFeatureProvider) SrsId() string {
	return provider.srsid
}
