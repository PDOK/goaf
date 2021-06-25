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
	Config         provider.Config
}

func NewGeopackageWithCommonProvider(api *openapi3.T, commonProvider provider.CommonProvider, crsMap map[string]string, config provider.Config) *GeoPackageProvider {
	gpkgFilePath := config.Datasource.Geopackage.File
	featureIdKey := config.Datasource.Geopackage.Fid

	return &GeoPackageProvider{
		CommonProvider: commonProvider,
		FilePath:       gpkgFilePath,
		CrsMap:         crsMap,
		FeatureIdKey:   featureIdKey,
		Api:            api,
		Config:         config,
	}
}

func (gp *GeoPackageProvider) Init() (err error) {
	gp.GeoPackage, err = NewGeoPackage(gp.FilePath, gp.FeatureIdKey)
	return
}
