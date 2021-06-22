package provider_gpkg

import (
	"fmt"
	"net/http"
	cg "oaf-server/codegen"
	pc "oaf-server/provider_common"
)

type DescribeCollectionProvider struct {
	data cg.Collection
}

func (provider *GeoPackageProvider) NewDescribeCollectionProvider(r *http.Request) (cg.Provider, error) {
	path := r.URL.Path // collections/{{collectionId}}
	ct := r.Header.Get("Content-Type")

	collectionId, _ := cg.ParametersForDescribeCollection(r)

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

		cInfo := cg.Collection{
			Id:          cn.Identifier,
			Title:       cn.Identifier,
			Description: cn.Description,
			Crs:         crss,
			Links:       []cg.Link{},
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

func (provider *DescribeCollectionProvider) SrsId() string {
	return "n.a"
}
