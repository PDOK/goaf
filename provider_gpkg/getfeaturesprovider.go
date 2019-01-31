package provider_gpkg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	cg "wfs3_server/codegen"
)

type GetFeaturesProvider struct {
	data FeatureCollectionGeoJSON
}

func (provider *GeoPackageProvider) NewGetFeaturesProvider(r *http.Request) (cg.Provider, error) {
	collectionId, limit, bbox, time, offset := cg.ParametersForGetFeatures(r)

	limitParam := provider.parseLimit(limit)
	offsetParam := provider.parseUint(offset, 0)
	bboxParam := provider.parseBBox(bbox, provider.GeoPackage.DefaultBBox)

	if time != "" {
		log.Println("Time selection currently not implemented")
	}

	path := r.URL.Path // collections/{{collectionId}}/items
	ct := r.Header.Get("Content-Type")

	p := &GetFeaturesProvider{}

	if ct == "" {
		ct = cg.JSONContentType
	}

	for _, cn := range provider.GeoPackage.Layers {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		featureId := uint64(0)
		fcGeoJSON, err := provider.GeoPackage.GetFeatures(r.Context(), provider.GeoPackage.DB, cn, collectionId, offsetParam, limitParam, featureId, bboxParam)

		if err != nil {
			return nil, err
		}

		requestParams := r.URL.Query()

		if int64(offsetParam) >= fcGeoJSON.NumberMatched && fcGeoJSON.NumberMatched > 0 {
			offsetParam = uint64(fcGeoJSON.NumberMatched - 1)
		}

		if int64(offsetParam) < 0 {
			offsetParam = 0
		}

		requestParams.Set("offset", fmt.Sprintf("%d", int64(offsetParam)))
		requestParams.Set("limit", fmt.Sprintf("%d", int64(limitParam)))

		// create links
		hrefBase := fmt.Sprintf("%s%s", provider.ServerEndpoint, path) // /collections
		links, _ := provider.createLinks(hrefBase, "self", ct)
		_ = provider.procesLinksForParams(links, requestParams)

		// next => offsetParam + limitParam < numbersMatched
		if (int64(offsetParam) + int64(limitParam)) < fcGeoJSON.NumberMatched {
			ilinks, _ := provider.createLinks(hrefBase, "next", ct)
			requestParams.Set("offset", fmt.Sprintf("%d", int64(offsetParam)+int64(limitParam)))
			_ = provider.procesLinksForParams(ilinks, requestParams)

			links = append(links, ilinks...)
		}

		// prev => offsetParam + limitParam < numbersMatched
		if int64(offsetParam) > 0 {
			ilinks, _ := provider.createLinks(hrefBase, "prev", ct)
			newOffset := int64(offsetParam) - int64(limitParam)
			if newOffset < 0 {
				newOffset = 0
			}

			requestParams.Set("offset", fmt.Sprintf("%d", newOffset))
			_ = provider.procesLinksForParams(ilinks, requestParams)

			links = append(links, ilinks...)
		}

		fcGeoJSON.Links = links

		p.data = fcGeoJSON
		break
	}

	return p, nil
}

func (provider *GetFeaturesProvider) Provide() (interface{}, error) {
	return provider.data, nil
}

func (provider *GetFeaturesProvider) MarshalJSON(interface{}) ([]byte, error) {
	return json.Marshal(provider.data)
}
func (provider *GetFeaturesProvider) MarshalHTML(interface{}) ([]byte, error) {
	// todo create html template pdok
	return json.Marshal(provider.data)
}
