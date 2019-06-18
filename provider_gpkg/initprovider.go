package provider_gpkg

import (
	"log"
	"net/url"
	"strconv"
	"strings"
	cg "wfs3_server/codegen"

	_ "github.com/mattn/go-sqlite3"
)

type GeoPackageProvider struct {
	FilePath           string
	GeoPackage         GeoPackage
	FeatureIdKey       string
	CrsMap             map[string]string
	serviceEndpoint    string
	serviceSpecPath    string
	maxReturnLimit     uint64
	defaultReturnLimit uint64
}

func NewGeopackageProvider(serviceEndpoint, serviceSpecPath, gpkgFilePath string, crsMap map[string]string, featureIdKey string, defaultReturnLimit uint64, maxReturnLimit uint64) *GeoPackageProvider {
	return &GeoPackageProvider{
		FilePath:           gpkgFilePath,
		CrsMap:             crsMap,
		FeatureIdKey:       featureIdKey,
		serviceEndpoint:    serviceEndpoint,
		serviceSpecPath:    serviceSpecPath,
		defaultReturnLimit: defaultReturnLimit,
		maxReturnLimit:     maxReturnLimit,
	}
}

func (provider *GeoPackageProvider) Init() (err error) {
	provider.GeoPackage, err = NewGeoPackage(provider.FilePath, provider.FeatureIdKey)
	return
}

func (provider *GeoPackageProvider) procesLinksForParams(links []cg.Link, queryParams url.Values) error {
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

func (provider *GeoPackageProvider) parseLimit(limit string) uint64 {
	limitParam := provider.defaultReturnLimit
	if limit != "" {
		newValue, err := strconv.ParseInt(limit, 10, 64)
		if err == nil && uint64(newValue) < provider.maxReturnLimit {
			limitParam = uint64(newValue)
		} else {
			limitParam = provider.maxReturnLimit
		}
	}
	return limitParam
}

func (provider *GeoPackageProvider) parseUint(stringValue string, defaultValue uint64) uint64 {
	if stringValue == "" {
		return defaultValue
	}
	value, err := strconv.ParseUint(stringValue, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func (provider *GeoPackageProvider) parseFloat64(stringValue string, defaultValue float64) float64 {
	if stringValue == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func (provider *GeoPackageProvider) parseBBox(stringValue string, defaultValue []float64) []float64 {
	if stringValue == "" {
		return provider.GeoPackage.DefaultBBox
	}
	bboxValues := strings.Split(stringValue, ",")
	if len(bboxValues) != 4 {
		return provider.GeoPackage.DefaultBBox
	}

	value := make([]float64, len(bboxValues))
	for i, v := range bboxValues {
		value[i] = provider.parseFloat64(v, value[i])
	}

	return value
}

func (provider *GeoPackageProvider) createLinks(path, rel, ct string) ([]cg.Link, error) {

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

func (provider *GeoPackageProvider) ctLink(baselink, contentType string) (string, error) {

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

func (provider *GeoPackageProvider) contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
