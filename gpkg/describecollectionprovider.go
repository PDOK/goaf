package gpkg

import (
	"net/http"
	"oaf-server/codegen"
	"oaf-server/provider"
)

func (gp *GeoPackageProvider) NewDescribeCollectionProvider(r *http.Request) (codegen.Provider, error) {
	return provider.NewDescribeCollectionProvider(gp.Config)(r)
}
