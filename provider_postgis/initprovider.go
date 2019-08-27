package provider_postgis

import (
	_ "github.com/mattn/go-sqlite3"
	"wfs3_server/provider_common"
)

type PostgisProvider struct {
	CommonProvider provider_common.CommonProvider
	PostGis        Postgis
	CrsMap         map[string]string
	configFilePath string
	connectionStr  string
}

func NewPostgisWithCommonProvider(commonProvider provider_common.CommonProvider, configPath, connectionStr string) *PostgisProvider {
	return &PostgisProvider{
		CrsMap:         map[string]string{"4326": "http://wfww.opengis.net/def/crs/OGC/1.3/CRS84"},
		configFilePath: configPath,
		connectionStr:  connectionStr,
		CommonProvider: commonProvider,
	}
}

func (provider *PostgisProvider) Init() (err error) {
	provider.PostGis, err = NewPostgis(provider.configFilePath, provider.connectionStr)
	return
}
