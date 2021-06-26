package gpkg

import (
	"log"
	"oaf-server/provider"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/imdario/mergo"
)

type GeoPackageProvider struct {
	CommonProvider provider.CommonProvider
	GeoPackage     GeoPackage
	Config         provider.Config
	Api            *openapi3.T
	ApiProcessed   *openapi3.T
}

func NewGeopackageWithCommonProvider(api *openapi3.T, commonProvider provider.CommonProvider, config provider.Config) *GeoPackageProvider {
	return &GeoPackageProvider{
		CommonProvider: commonProvider,
		Api:            api,
		Config:         config,
	}
}

func (gp *GeoPackageProvider) Init() (err error) {
	gp.GeoPackage, err = NewGeoPackage(gp.Config.Datasource.Geopackage.File, gp.Config.Datasource.Geopackage.Fid)

	collections := gp.Config.Datasource.Collections

	if len(collections) != 0 {
		for _, gpkgc := range gp.GeoPackage.Collections {
			for _, configc := range collections {
				if gpkgc.Identifier == configc.Identifier {
					err = mergo.Merge(&configc, gpkgc)
					if err != nil {
						log.Fatalln(err)
					}
				}
			}
		}
	} else {
		collections = gp.GeoPackage.Collections
	}

	gp.ApiProcessed = provider.CreateProvidesSpecificParameters(gp.Api, &collections)
	return
}
