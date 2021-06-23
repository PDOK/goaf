package gpkg

import (
	"oaf-server/provider"

	"github.com/getkin/kin-openapi/openapi3"
)

type GeoPackageProvider struct {
	CommonProvider provider.CommonProvider
	FilePath       string
	GeoPackage     GeoPackage
	FeatureIdKey   string
	CrsMap         map[string]string
	Api            *openapi3.T
}

func NewGeopackageWithCommonProvider(api *openapi3.T, commonProvider provider.CommonProvider, gpkgFilePath string, crsMap map[string]string, featureIdKey string) *GeoPackageProvider {
	return &GeoPackageProvider{
		CommonProvider: commonProvider,
		FilePath:       gpkgFilePath,
		CrsMap:         crsMap,
		FeatureIdKey:   featureIdKey,
		Api:            api,
	}
}

func (gp *GeoPackageProvider) Init() (err error) {
	gp.GeoPackage, err = NewGeoPackage(gp.FilePath, gp.FeatureIdKey)
	return
}
