package gpkg

import (
	"oaf-server/provider"

	"github.com/getkin/kin-openapi/openapi3"
)

type GeoPackageProvider struct {
	CommonProvider provider.CommonProvider
	GeoPackage     GeoPackage
	CrsMap         map[string]string
	Api            *openapi3.T
	Config         provider.Config
}

func NewGeopackageWithCommonProvider(api *openapi3.T, commonProvider provider.CommonProvider, crsMap map[string]string, config provider.Config) *GeoPackageProvider {
	return &GeoPackageProvider{
		CommonProvider: commonProvider,
		CrsMap:         crsMap,
		Api:            api,
		Config:         config,
	}
}

func (gp *GeoPackageProvider) Init() (err error) {
	gp.GeoPackage, err = NewGeoPackage(gp.Config.Datasource.Geopackage.File, gp.Config.Datasource.Geopackage.Fid)
	return
}
