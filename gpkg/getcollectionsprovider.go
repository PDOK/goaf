package gpkg

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

func (gp *GeoPackageProvider) NewGetCollectionsProvider(r *http.Request) (codegen.Provider, error) {

	path := r.URL.Path // collections
	ct := r.Header.Get("Content-Type")

	p := &GetCollectionsProvider{}
	p.contenttype = ct

	csInfo := codegen.Collections{Links: []codegen.Link{}, Collections: []codegen.Collection{}}
	// create Links
	hrefBase := fmt.Sprintf("%s%s", gp.CommonProvider.ServiceEndpoint, path) // /collections
	links, _ := provider.CreateLinks("collections", hrefBase, "self", ct)
	csInfo.Links = append(csInfo.Links, links...)
	for _, cn := range gp.GeoPackage.Layers {
		clinks, _ := provider.CreateLinks("collection "+cn.Identifier, fmt.Sprintf("%s/%s", hrefBase, cn.Identifier), "item", ct)
		csInfo.Links = append(csInfo.Links, clinks...)
	}

	for _, cn := range gp.GeoPackage.Layers {

		crss := make([]string, 0)
		for _, v := range gp.CrsMap {
			crss = append(crss, v)
		}

		cInfo := codegen.Collection{
			Id:          cn.Identifier,
			Title:       cn.Identifier,
			Description: cn.Description,
			Crs:         crss,
			Links:       []codegen.Link{},
		}

		chrefBase := fmt.Sprintf("%s/%s", hrefBase, cn.Identifier)

		clinks, _ := provider.CreateLinks("collection "+cn.Identifier, chrefBase, "self", ct)
		cInfo.Links = append(cInfo.Links, clinks...)

		cihrefBase := fmt.Sprintf("%s/items", chrefBase)
		ilinks, _ := provider.CreateLinks("items "+cn.Identifier, cihrefBase, "item", ct)
		cInfo.Links = append(cInfo.Links, ilinks...)

		for _, c := range gp.Config.Datasource.Collections {
			if c.Identifier == cn.Identifier {
				if len(c.Links) != 0 {
					cInfo.Links = append(cInfo.Links, c.Links...)
				}
				break
			}
		}

		csInfo.Collections = append(csInfo.Collections, cInfo)
	}

	p.data = csInfo

	return p, nil
}

func (gp *GetCollectionsProvider) Provide() (interface{}, error) {
	return gp.data, nil
}

func (gp *GetCollectionsProvider) ContentType() string {
	return gp.contenttype
}

func (gp *GetCollectionsProvider) String() string {
	return "getcollections"
}

func (gp *GetCollectionsProvider) SrsId() string {
	return "n.a"
}
