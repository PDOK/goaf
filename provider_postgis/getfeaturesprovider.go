package provider_postgis

import (
	"fmt"
	"log"
	"net/http"
	cg "wfs3_server/codegen"
	pc "wfs3_server/provider_common"
)

type GetFeaturesProvider struct {
	data FeatureCollectionGeoJSON
}

func (provider *PostgisProvider) NewGetFeaturesProvider(r *http.Request) (cg.Provider, error) {
	collectionId, limit, bbox, time, offset := cg.ParametersForGetFeatures(r)

	limitParam := pc.ParseLimit(limit, provider.commonProvider.DefaultReturnLimit, provider.commonProvider.MaxReturnLimit)
	offsetParam := pc.ParseUint(offset, 0)
	bboxParam := pc.ParseBBox(bbox, provider.PostGis.BBox)

	if time != "" {
		log.Println("Time selection currently not implemented")
	}

	path := r.URL.Path // collections/{{collectionId}}/items
	ct := r.Header.Get("Content-Type")

	p := &GetFeaturesProvider{}

	for _, cn := range provider.PostGis.Layers {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		fcGeoJSON, err := provider.PostGis.GetFeatures(r.Context(), provider.PostGis.db, cn, collectionId, offsetParam, limitParam, nil, bboxParam)

		if err != nil {
			return nil, err
		}

		requestParams := r.URL.Query()

		if int64(offsetParam) < 0 {
			offsetParam = 0
		}

		requestParams.Set("offset", fmt.Sprintf("%d", int64(offsetParam)))
		requestParams.Set("limit", fmt.Sprintf("%d", int64(limitParam)))

		// create links
		hrefBase := fmt.Sprintf("%s%s", provider.commonProvider.ServiceEndpoint, path) // /collections

		links, _ := pc.CreateLinks(hrefBase, "self", ct)
		_ = pc.ProcesLinksForParams(links, requestParams)

		// next => offsetParam + limitParam < numbersMatched
		if (int64(limitParam)) == fcGeoJSON.NumberReturned {
			ilinks, _ := pc.CreateLinks(hrefBase, "next", ct)
			requestParams.Set("offset", fmt.Sprintf("%d", int64(offsetParam)+int64(limitParam)))
			_ = pc.ProcesLinksForParams(ilinks, requestParams)

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

func (provider *GetFeaturesProvider) String() string {
	return "getfeatures"
}