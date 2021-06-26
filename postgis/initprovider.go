package postgis

import (
	"oaf-server/core"

	"github.com/getkin/kin-openapi/openapi3"
	_ "github.com/mattn/go-sqlite3"
)

type PostgisProvider struct {
	PostGis      Postgis
	Config       core.Config
	Api          *openapi3.T
	ApiProcessed *openapi3.T
}

func NewPostgisWithCommonProvider(api *openapi3.T, config core.Config) *PostgisProvider {
	p := &PostgisProvider{
		Config: config,
		Api:    api,
	}
	return p
}

func (pg *PostgisProvider) Init() (err error) {
	pg.PostGis, err = NewPostgis(pg.Config)
	pg.ApiProcessed = core.CreateProvidesSpecificParameters(pg.Api, &pg.PostGis.Collections)
	return
}
