**PDOK server implementation of [OGCAPI-FEATURES](https://github.com/opengeospatial/WFS_FES/blob/master/core/examples/openapi/ogcapi-features-1-example1.yaml) EXAMPLE 1.**

A a GeoJSON implementation with a Geopackage as a data provider.

Inspiration en code copied from https://github.com/go-spatial/jivan and https://github.com/go-spatial/tegola

The main differences with regards to jivan is the data provider setup, some geopackage query speedups for larger Geopackages and
some tweaks for scanning the SQL features

The specification is a preliminary one, with `go generate` the routing based on api spec, provider interfaces en types structs and convenient parameter extractions are generated to stay easily up to date.

* FeatureCollectionGeoJSON is overridden in provider gpkg to use the github.com/go-spatial/geom/encoding/geojso equivalent for decoding blobs
* https://github.com/opengeospatial/WFS_FES/blob/master/core/openapi/ogcapi-features-1.yaml

example wfs-3.0 geopackage example: https://github.com/PDOK/wfs-3.0-gpkg

***Minimal config, gpkg tends to be relative small e.g. < 3 GB***

go run start.go -provider gpkg -gpkg tst/bgt_wgs84.gpkg

***More elaborate config optimised performance for huge db (10M+ records/collection)***

go run start.go -provider postgis -config config/config_postgis.yaml

parameters :
```
bindHost := flag.String("s", envString("BIND_HOST", "0.0.0.0"), "server internal bind address, default; 0.0.0.0")
bindPort := flag.Int("p", envInt("BIND_PORT",8080), "server internal bind address, default; 8080")

serviceEndpoint := flag.String("endpoint", envString("ENDPOINT","http://localhost:8080"), "server endpoint for proxy reasons, default; http://localhost:8080")
serviceSpecPath := flag.String("spec", envString("SERVICE_SPEC_PATH","spec/wfs3.0.yml"), "swagger openapi spec")
defaultReturnLimit := flag.Int("limit", envInt("LIMIT",100), "limit, default: 100")
maxReturnLimit := flag.Int("limitmax", envInt("LIMIT_MAX",500), "max limit, default: 1000")
providerName := flag.String("provider", envString("PROVIDER",""), "postgis or gpkg")
gpkgFilePath := flag.String("gpkg", envString("PATH_GPKG",""), "geopackage path")
crsMapFilePath := flag.String("crs", envString("PATH_CRS",""), "crs file path")
configFilePath := flag.String("config", envString("PATH_CONFIG",""), "configfile path")

featureIdKey := flag.String("featureId", envString("FEATURE_ID",""), "Default feature identification or else first column definition (fid)") //optional for gpkg provider 

```


