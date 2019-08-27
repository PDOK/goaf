package provider_postgis

import (
	"fmt"
	"net/http"
	. "wfs3_server/codegen"
	pc "wfs3_server/provider_common"
)

type DescribeCollectionProvider struct {
	data Collection
}

func (provider *PostgisProvider) NewDescribeCollectionProvider(r *http.Request) (Provider, error) {
	path := r.URL.Path // collections/{{collectionId}}
	ct := r.Header.Get("Content-Type")

	collectionId := ParametersForDescribeCollection(r)

	p := &DescribeCollectionProvider{}

	for _, cn := range provider.PostGis.Layers {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		cInfo := Collection{
			Id:          cn.Identifier,
			Title:       cn.Identifier,
			Description: cn.Description,
			Crs:         []string{},

			Links: []Link{},
		}

		// create links
		hrefBase := fmt.Sprintf("%s%s", provider.CommonProvider.ServiceEndpoint, path) // /collections
		links, _ := pc.CreateLinks(collectionId, hrefBase, "self", ct)

		cihrefBase := fmt.Sprintf("%s/items", hrefBase)
		ilinks, _ := pc.CreateLinks("items of "+collectionId, cihrefBase, "item", ct)
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