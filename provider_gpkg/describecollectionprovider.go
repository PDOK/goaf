package provider_gpkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	. "wfs3_server/codegen"
)

type DescribeCollectionProvider struct {
	data CollectionInfo
}

func (provider *GeoPackageProvider) NewDescribeCollectionProvider(r *http.Request) (Provider, error) {
	path := r.URL.Path // collections/{{collectionId}}
	ct := r.Header.Get("Content-Type")

	collectionId := ParametersForDescribeCollection(r)

	p := &DescribeCollectionProvider{}

	if ct == "" {
		ct = JSONContentType
	}
	for _, cn := range provider.GeoPackage.Layers {
		// maybe convert to map, but not thread safe!
		if cn.Identifier != collectionId {
			continue
		}

		cInfo := CollectionInfo{
			Name:        cn.Identifier,
			Title:       cn.Identifier,
			Description: cn.Description,
			Links:       []Link{},
		}

		// create links
		hrefBase := fmt.Sprintf("%s%s", provider.ServerEndpoint, path) // /collections
		links, _ := provider.createLinks(hrefBase, "self", ct)

		cihrefBase := fmt.Sprintf("%s/items", hrefBase)
		ilinks, _ := provider.createLinks(cihrefBase, "item", ct)
		cInfo.Links = append(links, ilinks...)

		p.data = cInfo
		break
	}

	return p, nil
}

func (provider *DescribeCollectionProvider) Provide() (interface{}, error) {
	return provider.data, nil
}

func (provider *DescribeCollectionProvider) MarshalJSON(interface{}) ([]byte, error) {
	return json.Marshal(provider.data)
}
func (provider *DescribeCollectionProvider) MarshalHTML(interface{}) ([]byte, error) {
	// todo create html template pdok
	return json.Marshal(provider.data)
}
