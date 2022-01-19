package graphql

import (
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
func (gp *GraphqlProvider) NewGetFeatureProvider(r *http.Request) (codegen.Provider, error) {

	collectionId, featureId, _ := codegen.ParametersForGetFeature(r)

	featureIdParam := featureId
	bboxParam := gp.Graphql.DefaultBBox

	p := &GetFeatureProvider{srsid: fmt.Sprintf("EPSG:%d", gp.Graphql.Srid)}

	path := r.URL.Path

	ct, err := core.GetContentType(r, p.String())
	if err != nil {
		return nil, err
	}

	p.contenttype = ct

	for _, cn := range gp.Graphql.Collections {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		fcGeoJSON, err := gp.Graphql.GetFeatures(r.Context(), gp.Graphql.Url, cn, collectionId, 0, 1, featureIdParam, bboxParam)

		if err != nil {
			return nil, err
		}

		if len(fcGeoJSON.Features) == 1 {

			feature := fcGeoJSON.Features[0]

			hrefBase := fmt.Sprintf("%s%s", gp.Config.Service.Url, path) // /collections
			links, _ := core.CreateFeatureLinks("feature", hrefBase, "self", ct)
			feature.Links = links

			p.data = feature
		}

		break
	}

	return p, nil
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
