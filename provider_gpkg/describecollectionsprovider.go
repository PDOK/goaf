package provider_gpkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	. "wfs3_server/codegen"
)

type DescribeCollectionsProvider struct {
	data Content
}

func (provider *GeoPackageProvider) NewDescribeCollectionsProvider(r *http.Request) (Provider, error) {

	path := r.URL.Path // collections
	ct := r.Header.Get("Content-Type")

	if ct == "" {
		ct = JSONContentType
	}

	p := &DescribeCollectionsProvider{}

	csInfo := Content{Links: []Link{}, Collections: []CollectionInfo{}}
	// create Links
	hrefBase := fmt.Sprintf("%s%s", provider.serviceEndpoint, path) // /collections
	links, _ := provider.createLinks(hrefBase, "self", ct)
	csInfo.Links = append(csInfo.Links, links...)
	for _, cn := range provider.GeoPackage.Layers {
		clinks, _ := provider.createLinks(fmt.Sprintf("%s/%s", hrefBase, cn.Identifier), "item", ct)
		csInfo.Links = append(csInfo.Links, clinks...)
	}

	for _, cn := range provider.GeoPackage.Layers {

		crss := make([]string, 0)
		for _, v := range provider.CrsMap {
			crss = append(crss, v)
		}

		cInfo := CollectionInfo{
			Name:        cn.Identifier,
			Title:       cn.Identifier,
			Description: cn.Description,
			Crs:         crss,
			Links:       []Link{},
		}

		chrefBase := fmt.Sprintf("%s/%s", hrefBase, cn.Identifier)

		clinks, _ := provider.createLinks(chrefBase, "self", ct)
		cInfo.Links = append(cInfo.Links, clinks...)

		cihrefBase := fmt.Sprintf("%s/items", chrefBase)
		ilinks, _ := provider.createLinks(cihrefBase, "item", ct)
		cInfo.Links = append(cInfo.Links, ilinks...)
		csInfo.Collections = append(csInfo.Collections, cInfo)
	}

	p.data = csInfo

	return p, nil
}

func (provider *DescribeCollectionsProvider) Provide() (interface{}, error) {
	return provider.data, nil
}

func (provider *DescribeCollectionsProvider) MarshalJSON(interface{}) ([]byte, error) {
	return json.Marshal(provider.data)
}

func (provider *DescribeCollectionsProvider) MarshalHTML(interface{}) ([]byte, error) {
	return json.Marshal(provider.data)
}
