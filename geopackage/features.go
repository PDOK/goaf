package geopackage

import (
	"fmt"
	"log"
	"net/http"
	"oaf-server/codegen"
	"oaf-server/core"
)

type GetFeaturesProvider struct {
	data        core.FeatureCollectionGeoJSON
	srsid       string
	contenttype string
}

func (gp *GeoPackageProvider) NewGetFeaturesProvider(r *http.Request) (codegen.Provider, error) {
	collectionId, limit, offset, _, bbox, time := codegen.ParametersForGetFeatures(r)

	limitParam := core.ParseLimit(limit, uint64(gp.Config.DefaultFeatureLimit), uint64(gp.Config.MaxFeatureLimit))
	offsetParam := core.ParseUint(offset, 0)
	bboxParam := core.ParseBBox(bbox, gp.GeoPackage.DefaultBBox)

	if time != "" {
		log.Println("Time selection currently not implemented")
	}

	path := r.URL.Path // collections/{{collectionId}}/items
	p := &GetFeaturesProvider{srsid: fmt.Sprintf("EPSG:%d", gp.GeoPackage.Srid)}

	ct, err := core.GetContentType(r, p.String())
	if err != nil {
		return nil, err
	}

	p.contenttype = ct

	for _, cn := range gp.GeoPackage.Collections {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		fcGeoJSON, err := gp.GeoPackage.GetFeatures(r.Context(), gp.GeoPackage.DB, cn, collectionId, offsetParam, limitParam, nil, bboxParam)

		if err != nil {
			return nil, err
		}

		for _, feature := range fcGeoJSON.Features {
			hrefBase := fmt.Sprintf("%s%s/%v", gp.Config.Service.Url, path, feature.ID) // /collections
			links, _ := core.CreateFeatureLinks("feature", hrefBase, "self", ct)
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
		hrefBase := fmt.Sprintf("%s%s", gp.Config.Service.Url, path) // /collections
		links, _ := core.CreateFeatureLinks("features "+cn.Identifier, hrefBase, "self", ct)
		_ = core.ProcesLinksForParams(links, requestParams)

		// next => offsetParam + limitParam < numbersMatched
		if (int64(limitParam)) == fcGeoJSON.NumberReturned {
			ilinks, _ := core.CreateFeatureLinks("features "+cn.Identifier, hrefBase, "next", ct)
			requestParams.Set("offset", fmt.Sprintf("%d", int64(offsetParam)+int64(limitParam)))
			_ = core.ProcesLinksForParams(ilinks, requestParams)

			links = append(links, ilinks...)
		}

		fcGeoJSON.Links = links
		fcGeoJSON.Crs = gp.Config.Crs[fmt.Sprintf("%d", cn.Srid)]

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
	return "features"
}

func (gfp *GetFeaturesProvider) SrsId() string {
	return gfp.srsid
}
