package geopackage

import (
	"log"
	"oaf-server/core"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/imdario/mergo"
)

type GeoPackageProvider struct {
	GeoPackage   GeoPackage
	Config       core.Config
	Api          *openapi3.T
	ApiProcessed *openapi3.T
}

func NewGeopackageWithCommonProvider(api *openapi3.T, config core.Config) *GeoPackageProvider {
	return &GeoPackageProvider{
		Api:    api,
		Config: config,
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

	gp.ApiProcessed = core.CreateProvidesSpecificParameters(gp.Api, &collections)
	return
}
