package postgis

import (
	"fmt"
	"net/http"
	"oaf-server/codegen"
	"oaf-server/provider"
)

type GetCollectionsProvider struct {
	data        codegen.Collections
	contenttype string
}

func (pp *PostgisProvider) NewGetCollectionsProvider(r *http.Request) (codegen.Provider, error) {

	path := r.URL.Path // collections
	ct := r.Header.Get("Content-Type")

	p := &GetCollectionsProvider{}
	p.contenttype = ct

	csInfo := codegen.Collections{Links: []codegen.Link{}, Collections: []codegen.Collection{}}
	// create Links
	hrefBase := fmt.Sprintf("%s%s", pp.CommonProvider.ServiceEndpoint, path) // /collections
	links, _ := provider.CreateLinks("collections ", hrefBase, "self", ct)
	csInfo.Links = append(csInfo.Links, links...)
	for _, cn := range pp.PostGis.Layers {
		clinks, _ := provider.CreateLinks("collection "+cn.Identifier, fmt.Sprintf("%s/%s", hrefBase, cn.Identifier), "item", ct)
		csInfo.Links = append(csInfo.Links, clinks...)
	}

	for _, cn := range pp.PostGis.Layers {

		cInfo := codegen.Collection{
			Id:          cn.Identifier,
			Title:       cn.Identifier,
			Description: cn.Description,
			Crs:         []string{},
			Links:       []codegen.Link{},
		}

		chrefBase := fmt.Sprintf("%s/%s", hrefBase, cn.Identifier)

		clinks, _ := provider.CreateLinks("collection "+cn.Identifier, chrefBase, "self", ct)
		cInfo.Links = append(cInfo.Links, clinks...)

		cihrefBase := fmt.Sprintf("%s/items", chrefBase)
		ilinks, _ := provider.CreateLinks("items "+cn.Identifier, cihrefBase, "item", ct)
		cInfo.Links = append(cInfo.Links, ilinks...)
		csInfo.Collections = append(csInfo.Collections, cInfo)
	}

	p.data = csInfo

	return p, nil
}

func (gcp *GetCollectionsProvider) Provide() (interface{}, error) {
	return gcp.data, nil
}

func (gcp *GetCollectionsProvider) ContentType() string {
	return gcp.contenttype
}

func (gcp *GetCollectionsProvider) String() string {
	return "getcollections"
}

func (gcp *GetCollectionsProvider) SrsId() string {
	return "n.a."
}
