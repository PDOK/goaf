package postgis

import (
	"oaf-server/provider"

	"github.com/getkin/kin-openapi/openapi3"
	_ "github.com/mattn/go-sqlite3"
)

type PostgisProvider struct {
	CommonProvider provider.CommonProvider
	PostGis        Postgis
	Config         provider.Config
	Api            *openapi3.T
	ApiProcessed   *openapi3.T
}

func NewPostgisWithCommonProvider(api *openapi3.T, commonProvider provider.CommonProvider, config provider.Config) *PostgisProvider {
	p := &PostgisProvider{
		Config:         config,
		CommonProvider: commonProvider,
		Api:            api,
	}
	return p
}

func (pg *PostgisProvider) Init() (err error) {
	pg.PostGis, err = NewPostgis(pg.Config)
	pg.ApiProcessed = provider.CreateProvidesSpecificParameters(pg.Api, &pg.PostGis.Collections)
	return
}
