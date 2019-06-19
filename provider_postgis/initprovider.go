package provider_postgis

import (
	_ "github.com/mattn/go-sqlite3"
)

type PostgisProvider struct {
	PostGis            Postgis
	CrsMap             map[string]string
	serviceSpecPath    string
	configFilePath     string
	serviceEndpoint    string
	maxReturnLimit     uint64
	defaultReturnLimit uint64
}

func NewPostgisProvider(serviceEndpoint, servicespecPath, configPath string, defaultReturnLimit uint64, maxReturnLimit uint64) *PostgisProvider {
	return &PostgisProvider{
		CrsMap:             map[string]string{"4326": "http://wfww.opengis.net/def/crs/OGC/1.3/CRS84"},
		configFilePath:     configPath,
		serviceEndpoint:    serviceEndpoint,
		serviceSpecPath:    servicespecPath,
		defaultReturnLimit: defaultReturnLimit,
		maxReturnLimit:     maxReturnLimit,
	}
}

func (provider *PostgisProvider) Init() (err error) {
	provider.PostGis, err = NewPostgis(provider.configFilePath)
	return
}
