package provider_postgis

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/url"
	"strconv"
	"strings"
	cg "wfs3_server/codegen"
)

type PostgisProvider struct {
	ServerEndpoint  string
	ServiceSpecPath string
	ConnectionStr   string
	FeatureTables   []string
	PostGis         Postgis
	MaxLimit        uint64
	DefaultLimit    uint64
}

func (provider *PostgisProvider) Init() (err error) {
	provider.PostGis, err = NewPostgis(provider.ConnectionStr, provider.FeatureTables)
	return
}

func (provider *PostgisProvider) procesLinksForParams(links []cg.Link, queryParams url.Values) error {
	for l := range links {
		spath, err := url.Parse(links[l].Href)
		if err != nil {
			return err
		}
		values := spath.Query()

		for k, v := range queryParams {
			if k == "f" {
				continue
			}
			values.Add(k, v[0])
		}
		spath.RawQuery = values.Encode()
		links[l].Href = spath.String()
	}

	return nil

}

func (provider *PostgisProvider) parseLimit(limit string) uint64 {
	limitParam := provider.DefaultLimit
	if limit != "" {
		newValue, err := strconv.ParseInt(limit, 10, 64)
		if err == nil && uint64(newValue) < provider.MaxLimit {
			limitParam = uint64(newValue)
		} else {
			limitParam = provider.MaxLimit
		}
	}
	return limitParam
}

func (provider *PostgisProvider) parseUint(stringValue string, defaultValue uint64) uint64 {
	if stringValue == "" {
		return defaultValue
	}
	value, err := strconv.ParseUint(stringValue, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func (provider *PostgisProvider) parseFloat64(stringValue string, defaultValue float64) float64 {
	if stringValue == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func (provider *PostgisProvider) parseBBox(stringValue string, defaultValue []float64) []float64 {
	if stringValue == "" {
		return provider.PostGis.DefaultBBox
	}
	bboxValues := strings.Split(stringValue, ",")
	if len(bboxValues) != 4 {
		return provider.PostGis.DefaultBBox
	}

	value := make([]float64, len(bboxValues))
	for i, v := range bboxValues {
		value[i] = provider.parseFloat64(v, value[i])
	}

	return value
}

func (provider *PostgisProvider) createLinks(path, rel, ct string) ([]cg.Link, error) {

	links := make([]cg.Link, 0)

	links = append(links, cg.Link{Rel: rel, Href: path, Type: ct})

	if rel == "self" {
		rel = "alternate"
	}

	if rel != "self" {
		return links, nil
	}

	for _, sct := range cg.SupportedContentTypes {
		if ct == sct {
			continue
		}
		href, err := provider.ctLink(path, sct)
		if err != nil {
			return nil, err
		}

		links = append(links, cg.Link{Rel: rel, Href: href, Type: sct})
	}

	return links, nil
}

func (provider *PostgisProvider) ctLink(baselink, contentType string) (string, error) {

	u, err := url.Parse(baselink)
	if err != nil {
		log.Printf("Invalid link '%v', will return empty string.", baselink)
		return "", err
	}
	q := u.Query()

	var l string
	switch contentType {
	case cg.JSONContentType:
	default:
		q["f"] = []string{contentType}
	}

	u.RawQuery = q.Encode()
	l = u.String()
	return l, nil
}

func (provider *PostgisProvider) contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
