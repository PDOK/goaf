package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"oaf-server/core"
	"strings"
)

// Graphql configuration
type Graphql struct {
	ApplicationId string
	UserVersion   string
	Url           string
	Collections   []core.Collection
	BBox          [4]float64
	DefaultBBox   [4]float64
	Srid          int64
}

// NewGraphql returns the Graphql build on the given Config
func NewGraphql(config core.Config) (Graphql, error) {

	graphql := Graphql{}

	graphql.ApplicationId = config.ApplicationId
	graphql.UserVersion = config.UserVersion

	graphql.Collections = config.Datasource.Collections
	graphql.BBox = config.Datasource.BBox
	graphql.Srid = int64(config.Datasource.Srid)

	graphql.Url = config.Datasource.Graphql.URL

	return graphql, nil
}

func (graphql *Graphql) GetCollections(url string) (result []core.Collection, err error) {

	if graphql.Collections != nil {
		result = graphql.Collections
		err = nil
		return
	}

	// Shouldn't get here
	// needs to be configured through config file
	return
}

// GetFeatures return the FeatureCollection
func (graphql Graphql) GetFeatures(ctx context.Context, url string, collection core.Collection, collectionId string, offset uint64, limit uint64, featureId interface{}, bbox [4]float64) (result *core.FeatureCollection, err error) {

	if featureId != nil {
		return getfeature(url, fmt.Sprintf(`%s`, featureId))
	} else {
		return getfeatures(url, offset, limit, bbox)
	}
}

func getfeature(url, featureId string) (result *core.FeatureCollection, err error) {

	body := strings.NewReader(
		fmt.Sprintf(
			`{
				gebouw(identificatie: "%s") {
					identificatie
					status
					oorspronkelijkBouwjaar
					geometrie (srid: 9067){
						asWKB
					}
					geregistreerdMet{
						beginGeldigheid
					}
				}
			 }`, featureId),
	)

	resp, err := http.Post(url, `application/json`, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	mapresponseonfeatures(resp.Body)

	return nil, err
}

func getfeatures(url string, offset, limit uint64, bbox [4]float64) (result *core.FeatureCollection, err error) {
	return nil, nil
}

func mapresponseonfeatures(response io.ReadCloser) (result *core.FeatureCollection, err error) {

	type Gebouw struct {
		Identificatie          string `json:"identificatie"`
		Status                 string `json:"status"`
		OorspronkelijkBouwjaar string `json:"oorspronkelijkBouwjaar"`
		Geometrie              struct {
			AsWKB string `json:"asWKB"`
		} `json:"geometrie"`
		GeregistreerdMet struct {
			BeginGeldigheid string `json:"beginGeldigheid"`
		} `json:"geregistreerdMet"`
	}

	type Data struct {
		GebouwCollectie struct {
			Nodes []Gebouw `json:"nodes"`
		} `json:"gebouwCollectie"`
		Gebouw `json:"gebouw"`
	}

	data := Data{}
	body, err := ioutil.ReadAll(response)
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(body), &data)

	return nil, nil
}
