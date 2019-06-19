
use the folowing config file structure :

```
{
    "ConnectionStr": "host=127.0.0.1 port=5432 password=bgt dbname=bgt sslmode=disable",
    "PostGis": {
        "ApplicationId": "POSTGIS",
        "Layers": [
            {
                "SchemaName": "latest",
                "TableName": "pand",
                "Identifier": "pand_geometrie_vlak", // display identifier of the collection
                "Description": "Nice description maybe a link to external page",
                "GeometryColumn": "bbox", // should be a simple geometry (geojson)
                "FeatureIDColumn": "ogc_fid", // unique identifier of the feature {itemId} *
                "BBoxGeometryColumn": "bbox", // simpel 2d envelope
                "OffsetColumn": "offset_id",
                "BBox": [
                    13603.33,
                    306900.151,
                    277924.306,
                    617112.488
                ],
                "SrsId": 28992,
                "Features": [
                    "column_names"
                ]
            }
        ],
        
        // encompassing bbox of all layers 
        "BBox": [
            13603.33,
            306900.151,
            277924.306,
            617112.488
        ],
        "SrsId": 28992
    },
    "CrsMap": {
        "4326": "http://wfww.opengis.net/def/crs/OGC/1.3/CRS84"
    }
}

```
ALTER TABLE .... ADD COLUMN geom geometry(Polygon,28992);
UPDATE .... SET geom=st_forcesfs(geometrie_vlak)

ALTER TABLE .... ADD COLUMN bbox geometry(Polygon,28992);
UPDATE .... SET bbox=ST_Envelope(ST_Force2D(geometrie_vlak))


go run start.go -provider postgis -config config/config_postgis.yaml