package provider_gpkg

import (
	_ "github.com/mattn/go-sqlite3"
)

type GeoPackageProvider struct {
	FilePath           string
	GeoPackage         GeoPackage
	FeatureIdKey       string
	CrsMap             map[string]string
	serviceEndpoint    string
	serviceSpecPath    string
	maxReturnLimit     uint64
	defaultReturnLimit uint64
}

func NewGeopackageProvider(serviceEndpoint, serviceSpecPath, gpkgFilePath string, crsMap map[string]string, featureIdKey string, defaultReturnLimit uint64, maxReturnLimit uint64) *GeoPackageProvider {
	return &GeoPackageProvider{
		FilePath:           gpkgFilePath,
		CrsMap:             crsMap,
		FeatureIdKey:       featureIdKey,
		serviceEndpoint:    serviceEndpoint,
		serviceSpecPath:    serviceSpecPath,
		defaultReturnLimit: defaultReturnLimit,
		maxReturnLimit:     maxReturnLimit,
	}
}

func (provider *GeoPackageProvider) Init() (err error) {
	provider.GeoPackage, err = NewGeoPackage(provider.FilePath, provider.FeatureIdKey)
	return
}
