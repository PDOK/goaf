package geopackage

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
	p := &GetCollectionsProvider{}

	ct, err := provider.GetContentType(r, p.String())
	if err != nil {
		return nil, err
	}

	p.contenttype = ct

	csInfo := codegen.Collections{Links: []codegen.Link{}, Collections: []codegen.Collection{}}
	// create Links
	hrefBase := fmt.Sprintf("%s%s", gp.Config.Service.Url, path) // /collections
	links, _ := provider.CreateLinks("collections", p.String(), hrefBase, "self", ct)
	csInfo.Links = append(csInfo.Links, links...)
	for _, cn := range gp.GeoPackage.Collections {
		clinks, _ := provider.CreateLinks("collection "+cn.Identifier, p.String(), fmt.Sprintf("%s/%s", hrefBase, cn.Identifier), "item", ct)
		csInfo.Links = append(csInfo.Links, clinks...)
	}

	for _, cn := range gp.GeoPackage.Collections {

		cInfo := codegen.Collection{
			Id:          cn.Identifier,
			Title:       cn.Identifier,
			Description: cn.Description,
			Crs:         []string{gp.Config.Crs[fmt.Sprintf("%d", cn.Srid)]},
			Links:       []codegen.Link{},
		}

		chrefBase := fmt.Sprintf("%s/%s", hrefBase, cn.Identifier)

		clinks, _ := provider.CreateLinks("collection "+cn.Identifier, p.String(), chrefBase, "self", ct)
		cInfo.Links = append(cInfo.Links, clinks...)

		cihrefBase := fmt.Sprintf("%s/items", chrefBase)
		ilinks, _ := provider.CreateLinks("items "+cn.Identifier, p.String(), cihrefBase, "item", ct)
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