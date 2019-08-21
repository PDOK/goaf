package provider_gpkg

import (
	"wfs3_server/provider_common"
)

type GeoPackageProvider struct {
	CommonProvider provider_common.CommonProvider
	FilePath       string
	GeoPackage     GeoPackage
	FeatureIdKey   string
	CrsMap         map[string]string
}

func NewGeopackageWithCommonProvider(commonProvider provider_common.CommonProvider, gpkgFilePath string, crsMap map[string]string, featureIdKey string) *GeoPackageProvider {
	return &GeoPackageProvider{
		CommonProvider: commonProvider,
		FilePath:       gpkgFilePath,
		CrsMap:         crsMap,
		FeatureIdKey:   featureIdKey,
	}
}

func (provider *GeoPackageProvider) Init() (err error) {
	provider.GeoPackage, err = NewGeoPackage(provider.FilePath, provider.FeatureIdKey)
	return
}
