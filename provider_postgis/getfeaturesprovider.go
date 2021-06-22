package provider_postgis

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	cg "oaf-server/codegen"
	pc "oaf-server/provider_common"
)

type GetFeaturesProvider struct {
	data  FeatureCollectionGeoJSON
	srsid string
}

func (provider *PostgisProvider) NewGetFeaturesProvider(r *http.Request) (cg.Provider, error) {

	collectionId, limit, offset, _, bbox, time := cg.ParametersForGetFeatures(r)

	limitParam := pc.ParseLimit(limit, provider.CommonProvider.DefaultReturnLimit, provider.CommonProvider.MaxReturnLimit)
	offsetParam := pc.ParseUint(offset, 0)
	bboxParam := pc.ParseBBox(bbox, provider.PostGis.BBox)

	if time != "" {
		log.Println("Time selection currently not implemented")
	}

	path := r.URL.Path // collections/{collectionId}/items
	ct := r.Header.Get("Content-Type")

	p := &GetFeaturesProvider{srsid: fmt.Sprintf("EPSG:%d", provider.PostGis.SrsId)}

	pathItem := provider.ApiProcessed.Paths.Find(path)
	if pathItem == nil {
		return p, errors.New("Invalid path :" + path)
	}

	for k := range r.URL.Query() {
		if notfound := pathItem.Get.Parameters.GetByInAndName("query", k) == nil; notfound {
			return p, errors.New("Invalid query parameter :" + k)
		}
	}

	for _, cn := range provider.PostGis.Layers {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		whereMap := make(map[string]string)
		for i := range cn.VendorSpecificParameters {
			if qpv, exists := r.URL.Query()[cn.VendorSpecificParameters[i]]; exists {
				whereMap[cn.VendorSpecificParameters[i]] = qpv[0]
			}
		}

		fcGeoJSON, err := provider.PostGis.GetFeatures(r.Context(), provider.PostGis.db, cn, whereMap, offsetParam, limitParam, nil, bboxParam)

		if err != nil {
			return nil, err
		}

		for _, feature := range fcGeoJSON.Features {
			hrefBase := fmt.Sprintf("%s%s/%v", provider.CommonProvider.ServiceEndpoint, path, feature.ID) // /collections
			links, _ := pc.CreateLinks("feature", hrefBase, "self", ct)
			feature.Links = links
		}

		requestParams := r.URL.Query()

		if int64(offsetParam) < 0 {
			offsetParam = 0
		}

		requestParams.Set("offset", fmt.Sprintf("%d", int64(offsetParam)))
		requestParams.Set("limit", fmt.Sprintf("%d", int64(limitParam)))

		// create links
		hrefBase := fmt.Sprintf("%s%s", provider.CommonProvider.ServiceEndpoint, path) // /collections

		links, _ := pc.CreateLinks("features "+cn.Identifier, hrefBase, "self", ct)
		_ = pc.ProcesLinksForParams(links, requestParams)

		// next => offsetParam + limitParam < numbersMatched
		if (int64(limitParam)) == fcGeoJSON.NumberReturned {

			ilinks, _ := pc.CreateLinks("next features "+cn.Identifier, hrefBase, "next", ct)
			requestParams.Set("offset", fmt.Sprintf("%d", fcGeoJSON.Offset))
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

func (provider *GetFeaturesProvider) SrsId() string {
	return provider.srsid
}
