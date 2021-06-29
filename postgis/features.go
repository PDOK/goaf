package postgis

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"oaf-server/codegen"
	"oaf-server/core"
)

// GetFeaturesProvider is returned by the NewGetFeaturesProvider
// containing the data, srsid and contenttype for the response
type GetFeaturesProvider struct {
	data        core.FeatureCollection
	srsid       string
	contenttype string
}

// NewGetFeaturesProvider handles the request and return the GetFeaturesProvider
func (pp *PostgisProvider) NewGetFeaturesProvider(r *http.Request) (codegen.Provider, error) {

	collectionId, limit, offset, _, bbox, time := codegen.ParametersForGetFeatures(r)

	limitParam := core.ParseLimit(limit, uint64(pp.Config.DefaultFeatureLimit), uint64(pp.Config.MaxFeatureLimit))
	offsetParam := core.ParseUint(offset, 0)
	bboxParam := core.ParseBBox(bbox, pp.PostGis.BBox)

	if time != "" {
		log.Println("Time selection currently not implemented")
	}

	path := r.URL.Path // collections/{collectionId}/items

	p := &GetFeaturesProvider{srsid: fmt.Sprintf("EPSG:%d", pp.PostGis.Srid)}
	ct, err := core.GetContentType(r, p.String())

	if err != nil {
		return nil, err
	}

	p.contenttype = ct

	pathItem := pp.ApiProcessed.Paths.Find(path)
	if pathItem == nil {
		return p, errors.New("Invalid path :" + path)
	}

	for k := range r.URL.Query() {
		if notfound := pathItem.Get.Parameters.GetByInAndName("query", k) == nil; notfound {
			return p, errors.New("Invalid query parameter :" + k)
		}
	}

	for _, cn := range pp.PostGis.Collections {
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

		fcGeoJSON, err := pp.PostGis.GetFeatures(r.Context(), pp.PostGis.db, cn, whereMap, offsetParam, limitParam, nil, bboxParam)

		if err != nil {
			return nil, err
		}

		for _, feature := range fcGeoJSON.Features {
			hrefBase := fmt.Sprintf("%s%s/%v", pp.Config.Service.Url, path, feature.ID) // /collections
			links, _ := core.CreateFeatureLinks("feature", hrefBase, "self", ct)
			feature.Links = links
		}

		requestParams := r.URL.Query()

		if int64(offsetParam) < 0 {
			offsetParam = 0
		}

		requestParams.Set("offset", fmt.Sprintf("%d", int64(offsetParam)))
		requestParams.Set("limit", fmt.Sprintf("%d", int64(limitParam)))

		// create links
		hrefBase := fmt.Sprintf("%s%s", pp.Config.Service.Url, path) // /collections

		links, _ := core.CreateFeatureLinks("features "+cn.Identifier, hrefBase, "self", ct)
		_ = core.ProcesLinksForParams(links, requestParams)

		// next => offsetParam + limitParam < numbersMatched
		if (int64(limitParam)) == fcGeoJSON.NumberReturned {

			ilinks, _ := core.CreateFeatureLinks("next features "+cn.Identifier, hrefBase, "next", ct)
			requestParams.Set("offset", fmt.Sprintf("%d", fcGeoJSON.Offset))
			_ = core.ProcesLinksForParams(ilinks, requestParams)

			links = append(links, ilinks...)
		}

		fcGeoJSON.Links = links

		p.data = fcGeoJSON
		break
	}

	return p, nil
}

// Provide provides the data
func (gfp *GetFeaturesProvider) Provide() (interface{}, error) {
	return gfp.data, nil
}

// ContentType returns the ContentType
func (gfp *GetFeaturesProvider) ContentType() string {
	return gfp.contenttype
}

// String returns the provider name
func (gfp *GetFeaturesProvider) String() string {
	return "features"
}

// SrsId returns the srsid
func (gfp *GetFeaturesProvider) SrsId() string {
	return gfp.srsid
}
