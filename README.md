# GOAF

[![GitHub license](https://img.shields.io/github/license/PDOK/goaf)](https://github.com/PDOK/goaf/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/PDOK/goaf.svg)](https://github.com/PDOK/goaf/releases)
[![Go Report Card](https://goreportcard.com/badge/PDOK/goaf)](https://goreportcard.com/report/PDOK/goaf)
[![Docker Pulls](https://img.shields.io/docker/pulls/pdok/goaf.svg)](https://hub.docker.com/r/pdok/goaf)

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
docker run --rm -v `pwd`/example:/example -c /example/config-addresses-gpkg-minimal.yaml -p 8080:8080 pdok/goaf:latest
```

## PostGis

More elaborate config optimised performance for huge db (10M+ features/collection)

```docker
docker run -v `pwd`/example:/example -c /example/config-addresses-postgis-localhost.yaml' -p 8080:8080 pdok/goaf:latest
```

### Example table

```sql
CREATE TABLE addresses.addresses
(
    fid text COLLATE pg_catalog."default" NOT NULL,
    offsetid bigint NOT NULL,
    properties jsonb,
    geom geometry,
    bbox geometry,
      
    CONSTRAINT addresses_addresses_pk PRIMARY KEY (fid)
)
WITH (
    OIDS = FALSE
)

CREATE INDEX addresses_geom_sidx ON addresses.addresses USING GIST (geom);
CREATE INDEX addresses_offsetid_idx ON addresses.addresses(offsetid);
```

## Generate

Some of the code is generated based on the given oas.yaml:

```bash
codegen/provider.go
codegen/types.go
server/routing.gen.go
```

```go
go generate generate/gen.go
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
