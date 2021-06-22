# PostGIS provider

```go
go run start.go -provider postgis -config example/config_postgis.yaml
```

use the folowing config file structure :

```yaml
ApplicationId: POSTGIS
ConnectionStr: 'host=127.0.0.1 port=5432 dbname=bgt sslmode=disable'
Layers:
  - SchemaName: latest #database schema name
    TableName: pand   #database table/view name
    Identifier: pand  #collection name in api response
    Description: BGT Swifterband #Description of the collection
    GeometryColumn: geom  #column containing the simple feature geometry
    GeometryType: GEOMETRY # has currently no use
    FeatureIDColumn: ogc_fid #the unique indexed identifier for a given feature
    BBoxGeometryColumn: bbox #extra column with boundingbox selection index for intersects query
    OffsetColumn: ogc_fid # extra column to determine next keyset paging, should be numeric, indexed and unique, could be equal to feature id
    BBox: [13603.33,306900.151,277924.306,617112.488] # Bounding box of all features can be used to display subset of features
    SrsId: 28992 #the projection of the geometry in db's
    # Features are column names which should be exposed in properties par of the reponse
    Features:
      - ogc_fid
      - gml_id
      - namespace
      - lokaalid
      - objectbegintijd
      - objecteindtijd
      - tijdstipregistratie
      - eindregistratie
      - lv_publicatiedatum
      - bronhouder
      - inonderzoek
      - relatievehoogteligging
      - bgt_status
      - plus_status
      - identificatiebagpnd
      - nummeraanduidingtekst
      - nummeraanduidinghoek
      - identificatiebagvbolaagstehuisnummer
      - identificatiebagvbohoogstehuisnummer
BBox: [13603.33,306900.151,277924.306,617112.488] # bounding box of all layers
SrsId: 28992 # and the corresponding projection

```sql
ALTER TABLE .... ADD COLUMN geom geometry(Polygon,28992);
UPDATE .... SET geom=st_forcesfs(geometrie_vlak)

ALTER TABLE .... ADD COLUMN bbox geometry(Polygon,28992);
UPDATE .... SET bbox=ST_Envelope(ST_Force2D(geometrie_vlak))
```
