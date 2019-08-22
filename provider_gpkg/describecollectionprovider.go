package provider_gpkg

import (
	"fmt"
	"net/http"
	. "wfs3_server/codegen"
	pc "wfs3_server/provider_common"
)

type DescribeCollectionProvider struct {
	data Collection
}

func (provider *GeoPackageProvider) NewDescribeCollectionProvider(r *http.Request) (Provider, error) {
	path := r.URL.Path // collections/{{collectionId}}
	ct := r.Header.Get("Content-Type")

	collectionId := ParametersForDescribeCollection(r)

	p := &DescribeCollectionProvider{}

	for _, cn := range provider.GeoPackage.Layers {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		crss := make([]string, 0)
		for _, v := range provider.CrsMap {
			crss = append(crss, v)
		}

		cInfo := Collection{
			Id:          cn.Identifier,
			Title:       cn.Identifier,
			Description: cn.Description,
			Crs:         crss,
			Links:       []Link{},
		}

		// create links
		hrefBase := fmt.Sprintf("%s%s", provider.CommonProvider.ServiceEndpoint, path) // /collections
		links, _ := pc.CreateLinks("collection "+cn.Identifier, hrefBase, "self", ct)

		cihrefBase := fmt.Sprintf("%s/items", hrefBase)
		ilinks, _ := pc.CreateLinks("items "+cn.Identifier, cihrefBase, "item", ct)
		cInfo.Links = append(links, ilinks...)

		p.data = cInfo
		break
	}

	return p, nil
}

func (provider *DescribeCollectionProvider) Provide() (interface{}, error) {
	return provider.data, nil
}

func (provider *DescribeCollectionProvider) String() string {
	return "describecollection"
}