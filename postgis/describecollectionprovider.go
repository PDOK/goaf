package postgis

import (
	"fmt"
	"net/http"
	"oaf-server/codegen"
	"oaf-server/provider"
)

type DescribeCollectionProvider struct {
	data        codegen.Collection
	contenttype string
}

func (pp *PostgisProvider) NewDescribeCollectionProvider(r *http.Request) (codegen.Provider, error) {
	path := r.URL.Path // collections/{{collectionId}}
	ct := r.Header.Get("Content-Type")

	collectionId, _ := codegen.ParametersForDescribeCollection(r)

	p := &DescribeCollectionProvider{}
	p.contenttype = ct

	for _, cn := range pp.PostGis.Collections {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		cInfo := codegen.Collection{
			Id:          cn.Identifier,
			Title:       cn.Identifier,
			Description: cn.Description,
			Crs:         []string{},

			Links: []codegen.Link{},
		}

		// create links
		hrefBase := fmt.Sprintf("%s%s", pp.CommonProvider.ServiceEndpoint, path) // /collections
		links, _ := provider.CreateLinks(collectionId, hrefBase, "self", ct)

		cihrefBase := fmt.Sprintf("%s/items", hrefBase)
		ilinks, _ := provider.CreateLinks("items of "+collectionId, cihrefBase, "item", ct)
		cInfo.Links = append(links, ilinks...)

		p.data = cInfo
		break
	}

	return p, nil
}

func (dcp *DescribeCollectionProvider) Provide() (interface{}, error) {
	return dcp.data, nil
}

func (dcp *DescribeCollectionProvider) ContentType() string {
	return dcp.contenttype
}

func (dcp *DescribeCollectionProvider) String() string {
	return "describecollection"
}

func (dcp *DescribeCollectionProvider) SrsId() string {
	return "n.a."
}
