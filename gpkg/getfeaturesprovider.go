package gpkg

import (
	"fmt"
	"log"
	"net/http"
	"oaf-server/codegen"
	"oaf-server/provider"
)

type GetFeaturesProvider struct {
	data        FeatureCollectionGeoJSON
	srsid       string
	contenttype string
}

func (gp *GeoPackageProvider) NewGetFeaturesProvider(r *http.Request) (codegen.Provider, error) {
	collectionId, limit, offset, _, bbox, time := codegen.ParametersForGetFeatures(r)

	limitParam := provider.ParseLimit(limit, gp.CommonProvider.DefaultReturnLimit, gp.CommonProvider.MaxReturnLimit)
	offsetParam := provider.ParseUint(offset, 0)
	bboxParam := provider.ParseBBox(bbox, gp.GeoPackage.DefaultBBox)

	if time != "" {
		log.Println("Time selection currently not implemented")
	}

	path := r.URL.Path // collections/{{collectionId}}/items
	ct := r.Header.Get("Content-Type")
	if ct == provider.JSONContentType {
		ct = provider.GEOJSONContentType
	}

	p := &GetFeaturesProvider{srsid: fmt.Sprintf("EPSG:%d", gp.GeoPackage.SrsId)}
	p.contenttype = ct

	for _, cn := range gp.GeoPackage.Layers {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		fcGeoJSON, err := gp.GeoPackage.GetFeatures(r.Context(), gp.GeoPackage.DB, cn, collectionId, offsetParam, limitParam, nil, bboxParam)

		if err != nil {
			return nil, err
		}

		for _, feature := range fcGeoJSON.Features {
			hrefBase := fmt.Sprintf("%s%s/%v", gp.CommonProvider.ServiceEndpoint, path, feature.ID) // /collections
			links, _ := provider.CreateFeatureLinks("feature", hrefBase, "self", ct)
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
		hrefBase := fmt.Sprintf("%s%s", gp.CommonProvider.ServiceEndpoint, path) // /collections
		links, _ := provider.CreateFeatureLinks("features "+cn.Identifier, hrefBase, "self", ct)
		_ = provider.ProcesLinksForParams(links, requestParams)

		// next => offsetParam + limitParam < numbersMatched
		if (int64(limitParam)) == fcGeoJSON.NumberReturned {
			ilinks, _ := provider.CreateFeatureLinks("features "+cn.Identifier, hrefBase, "next", ct)
			requestParams.Set("offset", fmt.Sprintf("%d", int64(offsetParam)+int64(limitParam)))
			_ = provider.ProcesLinksForParams(ilinks, requestParams)

			links = append(links, ilinks...)
		}

		fcGeoJSON.Links = links

		crsUri, ok := gp.CrsMap[fmt.Sprintf("%d", cn.SrsId)]
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

func (gfp *GetFeaturesProvider) Provide() (interface{}, error) {
	return gfp.data, nil
}

func (gfp *GetFeaturesProvider) ContentType() string {
	return gfp.contenttype
}

func (gfp *GetFeaturesProvider) String() string {
	return "getfeatures"
}

func (gfp *GetFeaturesProvider) SrsId() string {
	return gfp.srsid
}
