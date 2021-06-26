package core

import (
	"fmt"
	"net/http"
	"oaf-server/codegen"
)

type GetCollectionsProvider struct {
	data        codegen.Collections
	contenttype string
}

func NewGetCollectionsProvider(config Config) func(r *http.Request) (codegen.Provider, error) {

	return func(r *http.Request) (codegen.Provider, error) {
		path := r.URL.Path // collections

		p := &GetCollectionsProvider{}

		ct, err := GetContentType(r, p.String())
		if err != nil {
			return nil, err
		}

		p.contenttype = ct

		csInfo := codegen.Collections{Links: []codegen.Link{}, Collections: []codegen.Collection{}}
		// create Links
		hrefBase := fmt.Sprintf("%s%s", config.Service.Url, path) // /collections
		links, _ := CreateLinks("collections ", p.String(), hrefBase, "self", ct)
		csInfo.Links = append(csInfo.Links, links...)
		for _, cn := range config.Datasource.Collections {
			clinks, _ := CreateLinks("collection "+cn.Identifier, p.String(), fmt.Sprintf("%s/%s", hrefBase, cn.Identifier), "item", ct)
			csInfo.Links = append(csInfo.Links, clinks...)
		}

		for _, cn := range config.Datasource.Collections {

			cInfo := codegen.Collection{
				Id:          cn.Identifier,
				Title:       cn.Identifier,
				Description: cn.Description,
				Crs:         []string{config.Crs[fmt.Sprintf("%d", cn.Srid)]},
				Links:       []codegen.Link{},
			}

			chrefBase := fmt.Sprintf("%s/%s", hrefBase, cn.Identifier)

			clinks, _ := CreateLinks("collection "+cn.Identifier, p.String(), chrefBase, "self", ct)
			cInfo.Links = append(cInfo.Links, clinks...)

			cihrefBase := fmt.Sprintf("%s/items", chrefBase)
			ilinks, _ := CreateLinks("items "+cn.Identifier, p.String(), cihrefBase, "item", ct)
			cInfo.Links = append(cInfo.Links, ilinks...)

			for _, c := range config.Datasource.Collections {
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
