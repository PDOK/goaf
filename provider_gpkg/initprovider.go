package provider_gpkg

import (
	"oaf-server/provider_common"

	"github.com/getkin/kin-openapi/openapi3"
)

type GeoPackageProvider struct {
	CommonProvider provider_common.CommonProvider
	FilePath       string
	GeoPackage     GeoPackage
	FeatureIdKey   string
	CrsMap         map[string]string
	Api            *openapi3.T
}

func NewGeopackageWithCommonProvider(api *openapi3.T, commonProvider provider_common.CommonProvider, gpkgFilePath string, crsMap map[string]string, featureIdKey string) *GeoPackageProvider {
	return &GeoPackageProvider{
		CommonProvider: commonProvider,
		FilePath:       gpkgFilePath,
		CrsMap:         crsMap,
		FeatureIdKey:   featureIdKey,
		Api:            api,
	}
}

func (provider *GeoPackageProvider) Init() (err error) {
	provider.GeoPackage, err = NewGeoPackage(provider.FilePath, provider.FeatureIdKey)
	return
}
