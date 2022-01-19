package graphql

import (
	"log"
	"oaf-server/core"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/imdario/mergo"
	_ "github.com/mattn/go-sqlite3"
)

// GraphqlProvider
type GraphqlProvider struct {
	Graphql      Graphql
	Config       core.Config
	Api          *openapi3.T
	ApiProcessed *openapi3.T
}

// NewGraphqlWithCommonProvider returns a new PostgisProvider set with the
// given config and OAS3 spec
func NewGraphqlWithCommonProvider(api *openapi3.T, config core.Config) *GraphqlProvider {
	g := &GraphqlProvider{
		Config: config,
		Api:    api,
	}
	return g
}

// Init initialize the Graphql backend
// and processed the OAS3 spec with the available collections
func (g *GraphqlProvider) Init() (err error) {
	g.Graphql, err = NewGraphql(g.Config)

	collections := g.Config.Datasource.Collections

	if len(collections) != 0 {
		for _, gc := range g.Graphql.Collections {
			for _, configc := range collections {
				if gc.Identifier == configc.Identifier {
					err = mergo.Merge(&configc, gc)
					if err != nil {
						log.Fatalln(err)
					}
				}
			}
		}
		g.Config.Datasource.Collections = collections
	} else {
		g.Config.Datasource.Collections = g.Graphql.Collections
	}

	g.ApiProcessed = core.CreateProvidesSpecificParameters(g.Api, &g.Graphql.Collections)
	return
}
