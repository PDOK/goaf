package core

import (
	"fmt"
	"net/http"
	"oaf-server/codegen"
)

type DescribeCollectionProvider struct {
	data        codegen.Collection
	contenttype string
}

func NewDescribeCollectionProvider(config Config) func(r *http.Request) (codegen.Provider, error) {

	return func(r *http.Request) (codegen.Provider, error) {
		path := r.URL.Path // collections/{{collectionId}}

		collectionId, _ := codegen.ParametersForDescribeCollection(r)

		p := &DescribeCollectionProvider{}

		ct, err := GetContentType(r, p.String())
		if err != nil {
			return nil, err
		}
		p.contenttype = ct

		for _, cn := range config.Datasource.Collections {
			// maybe convert to map, but not thread safe!
			if cn.Identifier != collectionId {
				continue
			}

			cInfo := codegen.Collection{
				Id:          cn.Identifier,
				Title:       cn.Identifier,
				Description: cn.Description,
				Crs:         []string{config.Crs[fmt.Sprintf("%d", cn.Srid)]},
				Links:       []codegen.Link{},
			}

			// create links
			hrefBase := fmt.Sprintf("%s%s", config.Service.Url, path) // /collections
			links, _ := CreateLinks(collectionId, p.String(), hrefBase, "self", ct)

			cihrefBase := fmt.Sprintf("%s/items", hrefBase)
			ilinks, _ := CreateLinks("items of "+collectionId, p.String(), cihrefBase, "item", ct)
			cInfo.Links = append(links, ilinks...)

			for _, c := range config.Datasource.Collections {
				if c.Identifier == cn.Identifier {
					if len(c.Links) != 0 {
						cInfo.Links = append(cInfo.Links, c.Links...)
					}
					break
				}
			}

			p.data = cInfo
			break
		}

		return p, nil
	}

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
