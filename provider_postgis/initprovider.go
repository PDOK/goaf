package provider_postgis

import (
	"github.com/getkin/kin-openapi/openapi3"
	_ "github.com/mattn/go-sqlite3"
	"wfs3_server/provider_common"
)

type PostgisProvider struct {
	CommonProvider provider_common.CommonProvider
	PostGis        Postgis
	CrsMap         map[string]string
	configFilePath string
	connectionStr  string
	Api            *openapi3.Swagger
	ApiProcessed   *openapi3.Swagger
}

func NewPostgisWithCommonProvider(api *openapi3.Swagger, commonProvider provider_common.CommonProvider, configPath, connectionStr string) *PostgisProvider {
	p := &PostgisProvider{
		CrsMap:         map[string]string{"4326": "http://wfww.opengis.net/def/crs/OGC/1.3/CRS84"},
		configFilePath: configPath,
		connectionStr:  connectionStr,
		CommonProvider: commonProvider,
		Api:            api,
	}
	return p
}

func (provider *PostgisProvider) Init() (err error) {
	provider.PostGis, err = NewPostgis(provider.configFilePath, provider.connectionStr)
	provider.ApiProcessed = CreateProvidesSpecificParameters(provider)
	return
}
