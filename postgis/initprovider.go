package postgis

import (
	"oaf-server/provider"

	"github.com/getkin/kin-openapi/openapi3"
	_ "github.com/mattn/go-sqlite3"
)

type PostgisProvider struct {
	PostGis      Postgis
	Config       provider.Config
	Api          *openapi3.T
	ApiProcessed *openapi3.T
}

func NewPostgisWithCommonProvider(api *openapi3.T, config provider.Config) *PostgisProvider {
	p := &PostgisProvider{
		Config: config,
		Api:    api,
	}
	return p
}

func (pg *PostgisProvider) Init() (err error) {
	pg.PostGis, err = NewPostgis(pg.Config)
	pg.ApiProcessed = provider.CreateProvidesSpecificParameters(pg.Api, &pg.PostGis.Collections)
	return
}
