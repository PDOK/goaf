package provider_gpkg

import (
	"fmt"
	"log"
	"net/http"
	cg "wfs3_server/codegen"
	pc "wfs3_server/provider_common"
)

type GetFeaturesProvider struct {
	data  FeatureCollectionGeoJSON
	srsid string
}

func (provider *GeoPackageProvider) NewGetFeaturesProvider(r *http.Request) (cg.Provider, error) {
	collectionId, limit, offset, _, bbox, time := cg.ParametersForGetFeatures(r)

	limitParam := pc.ParseLimit(limit, provider.CommonProvider.DefaultReturnLimit, provider.CommonProvider.MaxReturnLimit)
	offsetParam := pc.ParseUint(offset, 0)
	bboxParam := pc.ParseBBox(bbox, provider.GeoPackage.DefaultBBox)

	if time != "" {
		log.Println("Time selection currently not implemented")
	}

	path := r.URL.Path // collections/{{collectionId}}/items
	ct := r.Header.Get("Content-Type")

	p := &GetFeaturesProvider{srsid: fmt.Sprintf("EPSG:%d", provider.GeoPackage.SrsId)}

	for _, cn := range provider.GeoPackage.Layers {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		fcGeoJSON, err := provider.GeoPackage.GetFeatures(r.Context(), provider.GeoPackage.DB, cn, collectionId, offsetParam, limitParam, nil, bboxParam)

		if err != nil {
			return nil, err
		}

		for _, feature := range fcGeoJSON.Features {
			hrefBase := fmt.Sprintf("%s%s/%v", provider.CommonProvider.ServiceEndpoint, path, feature.ID) // /collections
			links, _ := pc.CreateLinks("feature", hrefBase, "self", ct)
			feature.Links = links
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
		hrefBase := fmt.Sprintf("%s%s", provider.CommonProvider.ServiceEndpoint, path) // /collections
		links, _ := pc.CreateLinks("features "+cn.Identifier, hrefBase, "self", ct)
		_ = pc.ProcesLinksForParams(links, requestParams)

		// next => offsetParam + limitParam < numbersMatched
		if (int64(limitParam)) == fcGeoJSON.NumberReturned {
			ilinks, _ := pc.CreateLinks("features "+cn.Identifier, hrefBase, "next", ct)
			requestParams.Set("offset", fmt.Sprintf("%d", int64(offsetParam)+int64(limitParam)))
			_ = pc.ProcesLinksForParams(ilinks, requestParams)

			links = append(links, ilinks...)
		}

		fcGeoJSON.Links = links

		crsUri, ok := provider.CrsMap[fmt.Sprintf("%d", cn.SrsId)]
		if !ok {
			log.Printf("SRS ID: %s, not found", fmt.Sprintf("%d", cn.SrsId))
			crsUri = ""
		}
		fcGeoJSON.Crs = crsUri

		p.data = *fcGeoJSON
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
