# GOAF

[![GitHub license](https://img.shields.io/github/license/PDOK/goaf)](https://github.com/PDOK/goaf/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/PDOK/goaf.svg)](https://github.com/PDOK/goaf/releases)
[![Go Report Card](https://goreportcard.com/badge/PDOK/goaf)](https://goreportcard.com/report/PDOK/goaf)

Goaf is a [OGC API - Features](https://www.ogc.org/standards/ogcapi-features) implementation in golang.

## For who is it?

If you are looking for a very fast OGC API - Features application and are not afraid for some data tuning, this could be the application for you!

## Datasources

It supports the following datasources:

* [OGC GeoPackage](https://www.geopackage.org/)
* [Postgis](https://postgis.net/) (*Postgresql 9.4+)

## PDOK server implementation of [OGCAPI-FEATURES](https://github.com/opengeospatial/ogcapi-features/blob/master/core/examples/openapi/ogcapi-features-1-example1.yaml)

A a GeoJSON implementation with a Geopackage as a data provider.

The specification is a preliminary one, with `go generate` the routing based on api spec, provider interfaces en types structs and convenient parameter extractions are generated to stay easily up to date.

* FeatureCollectionGeoJSON is overridden in the GeoPackage provider to use the [geojson](https://github.com/go-spatial/geom/tree/master/encoding/geojson) equivalent for decoding blobs
* <https://github.com/opengeospatial/ogcapi-features/blob/master/core/openapi/ogcapi-features-1.yaml>

## Build

```docker
docker build -t pdok/goaf:latest .
```

## GeoPackage

The geopacakge provider is a minimal config for GeoPackages that tend to be relative small e.g. < 3 GB.

```docker
docker run --rm -v `pwd`/example:/example -e PROVIDER='gpkg' -e PATH_GPKG='/example/addresses.gpkg' -e ENDPOINT='http://localhost:8080' -p 8080:8080 pdok/goaf:latest
```

## PostGis

More elaborate config optimised performance for huge db (10M+ records/collection)

```docker
docker run -v `pwd`/example:/example -e CONNECTION='postgres://{user}:{password}@{host}:{port}/{database}?sslmode=disable' -e PROVIDER='postgis' -e PATH_CONFIG='/example/config_postgis.yaml' -e ENDPOINT='http://localhost:8080' -p 8080:8080 pdok/goaf:latest
```

example table

```sql
CREATE TABLE bgt_wfs3_v1.bak
(
    _id text COLLATE pg_catalog."default" NOT NULL,
    _version text COLLATE pg_catalog."default",
    properties jsonb,
    _geom geometry,
    _bbox geometry,
    _offset_id bigint NOT NULL DEFAULT nextval('bgt_wfs3_v1.bak__offset_id_seq'::regclass),
    _created timestamp without time zone,
    CONSTRAINT bak_pkey PRIMARY KEY (_id)
)
WITH (
    OIDS = FALSE
)
```

used parameters:

```go
bindHost := flag.String("s", envString("BIND_HOST", "0.0.0.0"), "server internal bind address, default; 0.0.0.0")
bindPort := flag.Int("p", envInt("BIND_PORT",8080), "server internal bind address, default; 8080")

serviceEndpoint := flag.String("endpoint", envString("ENDPOINT","http://localhost:8080"), "server endpoint for proxy reasons, default; http://localhost:8080")
serviceSpecPath := flag.String("spec", envString("SERVICE_SPEC_PATH","spec/oaf.yml"), "swagger openapi spec")
defaultReturnLimit := flag.Int("limit", envInt("LIMIT",100), "limit, default: 100")
maxReturnLimit := flag.Int("limitmax", envInt("LIMIT_MAX",500), "max limit, default: 1000")
providerName := flag.String("provider", envString("PROVIDER",""), "postgis or gpkg")
gpkgFilePath := flag.String("gpkg", envString("PATH_GPKG",""), "geopackage path")
crsMapFilePath := flag.String("crs", envString("PATH_CRS",""), "crs file path")
configFilePath := flag.String("config", envString("PATH_CONFIG",""), "configfile path")
connectionStr := flag.String("connection", envString("CONNECTION", ""), "configfile path")

featureIdKey := flag.String("featureId", envString("FEATURE_ID",""), "Default feature identification or else first column definition (fid)") //optional for gpkg provider 
```

## Test

```go
go test ./... -covermode=atomic
```

## How to Contribute

Make a pull request...

## License

Distributed under MIT License, please see license file within the code for more details.

## Thanks

Inspiration en code copied from:

* <https://github.com/go-spatial/jivan>
* <https://github.com/go-spatial/tegola>

The main differences with regards to jivan is the data provider setup, some geopackage query speedups for larger Geopackages and
some tweaks for scanning the SQL features
