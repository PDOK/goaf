package gpkg

import (
	"oaf-server/provider"

	"github.com/getkin/kin-openapi/openapi3"
)

type GeoPackageProvider struct {
	CommonProvider provider.CommonProvider
	GeoPackage     GeoPackage
	CrsMap         map[string]string
	Config         provider.Config
	Api            *openapi3.T
	ApiProcessed   *openapi3.T
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
	// gp.ApiProcessed = CreateProvidesSpecificParameters(gp.Api, gp.Config.Datasource.Collections)
	return
}
