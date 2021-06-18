package provider_postgis

import (
	"fmt"
	"net/http"
	cg "wfs3_server/codegen"
	pc "wfs3_server/provider_common"
)

type DescribeCollectionProvider struct {
	data cg.Collection
}

func (provider *PostgisProvider) NewDescribeCollectionProvider(r *http.Request) (cg.Provider, error) {
	path := r.URL.Path // collections/{{collectionId}}
	ct := r.Header.Get("Content-Type")

	collectionId, _ := cg.ParametersForDescribeCollection(r)

	p := &DescribeCollectionProvider{}

	for _, cn := range provider.PostGis.Layers {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		cInfo := cg.Collection{
			Id:          cn.Identifier,
			Title:       cn.Identifier,
			Description: cn.Description,
			Crs:         []string{},

			Links: []cg.Link{},
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

func (provider *DescribeCollectionProvider) SrsId() string {
	return "n.a."
}
