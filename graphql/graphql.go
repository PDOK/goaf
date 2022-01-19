package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"oaf-server/core"
	"strings"

	"github.com/go-spatial/geom/encoding/geojson"
	"github.com/go-spatial/geom/encoding/wkt"
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

	type graphqlreq struct {
		Query     string      `json:"query"`
		Variables interface{} `json:"variables"`
	}

	query := fmt.Sprintf(`{
            gebouw(identificatie: "%s") {
                identificatie
                status
                oorspronkelijkBouwjaar
                geometrie (srid: 9067) {
                    asWKT
                }
                geregistreerdMet{
                    beginGeldigheid
                }
            }
         }`, featureId)

	req := graphqlreq{Query: query, Variables: nil}

	b, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	resp, err := http.Post(url, `application/json`, bytes.NewReader(b))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return mapresponseonfeatures(body)
}

func getfeatures(url string, offset, limit uint64, bbox [4]float64) (result *core.FeatureCollection, err error) {

	type graphqlreq struct {
		Query     string      `json:"query"`
		Variables interface{} `json:"variables"`
	}

	query := fmt.Sprintf(`{
		    gebouwCollectie(first: %d, offset: %d) {
              nodes {
                identificatie
                status
                oorspronkelijkBouwjaar
                geometrie (srid: 9067) {
                    asWKT
                }
                geregistreerdMet{
                    beginGeldigheid
                }
			  }
            }
         }`, limit, offset)

	req := graphqlreq{Query: query, Variables: nil}

	b, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	resp, err := http.Post(url, `application/json`, bytes.NewReader(b))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Print(string(body))

	return mapresponseonfeatures(body)
}

func mapresponseonfeatures(body []byte) (result *core.FeatureCollection, err error) {

	type Gebouw struct {
		Identificatie          string `json:"identificatie"`
		Status                 string `json:"status"`
		OorspronkelijkBouwjaar string `json:"oorspronkelijkBouwjaar"`
		Geometrie              struct {
			AsWKT string `json:"asWKT"`
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

	type Body struct {
		Data Data `json:"data"`
	}

	data := Body{}
	json.Unmarshal([]byte(body), &data)

	var gebouwen []Gebouw

	if data.Data.GebouwCollectie.Nodes != nil {
		gebouwen = data.Data.GebouwCollectie.Nodes
	} else {
		gebouwen = []Gebouw{data.Data.Gebouw}
	}

	fs := []*core.Feature{}

	for _, gebouw := range gebouwen {

		geom, err := wkt.Decode(strings.NewReader(gebouw.Geometrie.AsWKT))
		if err != nil {
			return nil, err
		}

		prop := map[string]interface{}{
			"status":                 gebouw.Status,
			"oorspronkelijkbouwjaar": gebouw.OorspronkelijkBouwjaar,
		}

		f := core.Feature{
			ID: gebouw.Identificatie,
			Feature: geojson.Feature{
				Geometry:   geojson.Geometry{Geometry: geom},
				Properties: prop,
			},
		}

		fs = append(fs, &f)
	}

	fc := core.FeatureCollection{Features: fs}
	fc.NumberReturned = int64(len(fs))
	fc.Type = "FeatureCollection"

	return &fc, nil
}
