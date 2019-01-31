**PDOK server implementation of [OGC WFS 3.0](https://github.com/opengeospatial/WFS_FES).**

A a GeoJSON implementation with a Geopackage as a data provider.

Inspiration en code copied from https://github.com/go-spatial/jivan and https://github.com/go-spatial/tegola

The main differences with regards to jivan is the data provider setup, some geopackage query speedups for larger Geopackages and
some tweaks for scanning the SQL features

The specification is a preliminary one, with `go generate` the routing based on api spec, provider interfaces en types structs and convenient parameter extractions are generated to stay easily up to date.

* FeatureCollectionGeoJSON is overridden in provider gpkg to use the github.com/go-spatial/geom/encoding/geojso equivalent for decoding blobs
* https://raw.githubusercontent.com/opengeospatial/WFS_FES/master/openapi.yaml

**DOCKER**  

` docker build -t wfs . && docker run -p 8080:8080 wfs`

**example geopackage detail log**

2019/01/24 13:35:45 | GEOPACKAGE DETAILS 
2019/01/24 13:35:45 |
2019/01/24 13:35:45 |   FILE: /2019_gemeentegrenzen_kustlijn.gpkg, APPLICATION: GP10, VERSION: 0
2019/01/24 13:35:45 |
2019/01/24 13:35:45 |   NUMBER OF LAYERS: 1
2019/01/24 13:35:45 |
2019/01/24 13:35:45 |   LAYER: 1. ID: 2019_gemeentegrenzen_kustlijn, SRS_ID: 28992, TABLE: 2019_gemeentegrenzen_kustlijn PK: fid, FEATURES : [geom id gid code gemeentenaam]
2019/01/24 13:35:45 | 
2019/01/24 13:35:45 |   BBOX: [13565.400000,306846.000000,278026.000000,619233.000000], SRS_ID:28992
2019/01/24 13:35:45 |
2019/01/24 13:35:45 | SERVING ON: http://localhost:8080





