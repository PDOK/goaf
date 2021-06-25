package postgis

import (
	"oaf-server/provider"

	"github.com/getkin/kin-openapi/openapi3"
	_ "github.com/mattn/go-sqlite3"
)

type PostgisProvider struct {
	CommonProvider provider.CommonProvider
	PostGis        Postgis
	CrsMap         map[string]string
	Config         provider.Config
	Api            *openapi3.T
	ApiProcessed   *openapi3.T
}

func NewPostgisWithCommonProvider(api *openapi3.T, commonProvider provider.CommonProvider, config provider.Config) *PostgisProvider {
	p := &PostgisProvider{
		CrsMap:         map[string]string{"4326": "http://wfww.opengis.net/def/crs/OGC/1.3/CRS84"},
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
