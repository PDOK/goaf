package provider_common

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	cg "wfs3_server/codegen"
)

const (
	JSONContentType = "application/json"
	HTMLContentType = "text/html"
)

func GetContentTypes() map[string]string {
	ct := make(map[string]string)

	ct["json"] = JSONContentType
	ct["html"] = HTMLContentType

	return ct
}

func GetRelationMap() map[string]string {
	ct := make(map[string]string)

	ct["alternate"] = "Alternative"
	ct["self"] = "This"

	return ct
}

func GetContentFields() map[string]string {
	ct := make(map[string]string)

	ct[JSONContentType] = "json"
	ct[HTMLContentType] = "html"

	return ct
}


func ProcesLinksForParams(links []cg.Link, queryParams url.Values) error {
	for l := range links {
		path, err := url.Parse(links[l].Href)
		if err != nil {
			return err
		}
		values := path.Query()

		for k, v := range queryParams {
			if k == "f" {
				continue
			}
			values.Add(k, v[0])
		}
		path.RawQuery = values.Encode()
		links[l].Href = path.String()
	}

	return nil

}

func CreateLinks(title, hrefPath, rel, ct string) ([]cg.Link, error) {

	links := make([]cg.Link, 0)

	href, err := ctLink(hrefPath, GetContentFields()[ct])
	if err != nil {
		return links, err
	}
	links = append(links, cg.Link{Title: formatTitle(title, rel, GetContentFields()[ct]), Rel: rel, Href: href, Type: ct})

	if rel == "self" {
		rel = "alternate"
	}

	//if rel != "self" {
	//	return links, nil
	//}

	for k, sct := range GetContentTypes() {
		if ct == sct {
			continue
		}
		href, err := ctLink(hrefPath, k)
		if err != nil {
			return links, err
		}

		links = append(links, cg.Link{Title: formatTitle(title, rel, k), Rel: rel, Href: href, Type: sct})

		if rel == "self" {
			rel = "alternate"
			links = append(links, cg.Link{Title: formatTitle(title, rel, k), Rel: rel, Href: href, Type: sct})
		}
	}

	return links, nil
}

func formatTitle(title, rel, format string) string {
	relation := rel
	if "self" == rel {
		relation = "this"
	}
	return strings.ToLower(fmt.Sprintf("%s %s in %s format", relation, title, format))
}

func ctLink(baselink, contentType string) (string, error) {

	u, err := url.Parse(baselink)
	if err != nil {
		log.Printf("Invalid link '%v', will return empty string.", baselink)
		return "", err
	}
	q := u.Query()

	var l string
	switch contentType {
	default:
		q["f"] = []string{contentType}
	}

	u.RawQuery = q.Encode()
	l = u.String()
	return l, nil
}

// copied,tweaked from https://github.com/go-spatial/jivan
func ConvertFeatureID(v interface{}) (interface{}, error) {
	switch aval := v.(type) {
	case float64:
		return uint64(aval), nil
	case int64:
		return uint64(aval), nil
	case uint64:
		return aval, nil
	case uint:
		return uint64(aval), nil
	case int8:
		return uint64(aval), nil
	case uint8:
		return uint64(aval), nil
	case uint16:
		return uint64(aval), nil
	case int32:
		return uint64(aval), nil
	case uint32:
		return uint64(aval), nil
	case []byte:
		return string(aval), nil
	case string:
		return aval, nil

	default:
		return 0, errors.New(fmt.Sprintf("Cannot convert ID : %v", aval))
	}
}

func ParseLimit(limit string, defaultReturnLimit, maxReturnLimit uint64) uint64 {
	limitParam := defaultReturnLimit
	if limit != "" {
		newValue, err := strconv.ParseInt(limit, 10, 64)
		if err == nil && uint64(newValue) < maxReturnLimit {
			limitParam = uint64(newValue)
		} else {
			limitParam = maxReturnLimit
		}
	}
	return limitParam
}

func ParseUint(stringValue string, defaultValue uint64) uint64 {
	if stringValue == "" {
		return defaultValue
	}
	value, err := strconv.ParseUint(stringValue, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func ParseFloat64(stringValue string, defaultValue float64) float64 {
	if stringValue == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func ParseBBox(stringValue string, defaultValue []float64) []float64 {
	if stringValue == "" {
		return defaultValue
	}
	bboxValues := strings.Split(stringValue, ",")
	if len(bboxValues) != 4 {
		return defaultValue
	}

	value := make([]float64, len(bboxValues))
	for i, v := range bboxValues {
		value[i] = ParseFloat64(v, value[i])
	}

	return value
}

func UpperFirst(title string) string {
	return strings.Title(title)
}