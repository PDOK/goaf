package postgis

import (
	"errors"
	"fmt"
	"net/http"
	"oaf-server/codegen"
	"oaf-server/core"
)

// GetFeatureProvider is returned by the NewGetFeatureProvider
// containing the data, srsid and contenttype for the response
type GetFeatureProvider struct {
	data        *core.Feature
	srsid       string
	contenttype string
}

// NewGetFeatureProvider handles the request and return the GetFeatureProvider
func (pp *PostgisProvider) NewGetFeatureProvider(r *http.Request) (codegen.Provider, error) {

	collectionId, featureId, _ := codegen.ParametersForGetFeature(r)

	featureIdParam := featureId
	bboxParam := pp.PostGis.BBox

	p := &GetFeatureProvider{srsid: fmt.Sprintf("EPSG:%d", pp.PostGis.Srid)}

	path := r.URL.Path

	ct, err := core.GetContentType(r, p.String())
	if err != nil {
		return nil, err
	}

	p.contenttype = ct

	for _, cn := range pp.PostGis.Collections {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		pathItem := pp.ApiProcessed.Paths.Find("/collections/" + collectionId + "/items/{featureId}")
		if pathItem == nil {
			return p, errors.New("Invalid path :" + path)
		}

		for k := range r.URL.Query() {
			if notfound := pathItem.Get.Parameters.GetByInAndName("query", k) == nil; notfound {
				return p, errors.New("Invalid query parameter :" + k)
			}
		}

		whereMap := make(map[string]string)
		fcGeoJSON, err := pp.PostGis.GetFeatures(r.Context(), pp.PostGis.db, cn, whereMap, 0, 1, featureIdParam, bboxParam)

		if err != nil {
			return nil, err
		}

		if len(fcGeoJSON.Features) >= 1 {
			feature := fcGeoJSON.Features[0]

			hrefBase := fmt.Sprintf("%s%s", pp.Config.Service.Url, path) // /collections
			links, _ := core.CreateFeatureLinks("feature", hrefBase, "self", ct)
			feature.Links = links

			p.data = feature

		} else {
			return p, fmt.Errorf("feature with id: %s not found", string(featureId))
		}

		return p, nil
	}

	return p, errors.New("Cannot find collection : " + collectionId)
}

// Provide provides the data
func (gfp *GetFeatureProvider) Provide() (interface{}, error) {
	return gfp.data, nil
}

// ContentType returns the ContentType
func (gfp *GetFeatureProvider) ContentType() string {
	return gfp.contenttype
}

// String returns the provider name
func (gfp *GetFeatureProvider) String() string {
	return "feature"
}

// SrsId returns the srsid
func (gfp *GetFeatureProvider) SrsId() string {
	return gfp.srsid
}
