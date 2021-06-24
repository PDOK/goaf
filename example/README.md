# Example

## create schema

```sql
CREATE SCHEMA addresses;
```

## import data

```docker
docker run --rm -v `pwd`:/example osgeo/gdal:ubuntu-small-3.3.0 ogr2ogr -f PostgreSQL "PG:dbname=oaf host=host.docker.internal port=5432 user=docker password=docker SCHEMAS=addresses" -nln addresses_alternative_encoding -nlt POINT /example/addresses.gpkg
```

## Transform data

```sql
CREATE TABLE addresses.addresses AS
SELECT fid, 
       fid AS offsetid,
       json_build_object(
           'alternativeidentifier', alternativeidentifier,
           'validfrom', validfrom,
           'validto', validto,
           'beginlifespanversion', beginlifespanversion,
           'endlifespanversion', endlifespanversion,
           'building', building,
           'component_thoroughfarename', component_thoroughfarename,
           'component_postaldescriptor', component_postaldescriptor,
           'component_addressareaname', component_addressareaname,
           'component_adminunitname_1', component_adminunitname_1,
           'component_adminunitname_2', component_adminunitname_2,
           'component_adminunitname_3', component_adminunitname_3,
           'component_adminunitname_4', component_adminunitname_4,
           'component_adminunitname_5', component_adminunitname_5,
           'component_adminunitname_6', component_adminunitname_6,
           'locator_designator_addressnumber', locator_designator_addressnumber,
           'locator_designator_addressnumberextension', locator_designator_addressnumberextension,
           'locator_designator_addressnumber2ndextension', locator_designator_addressnumber2ndextension,
           'locator_level', locator_level,
           'locator_href', locator_href,
           'locator_designator_buildingidentifier', locator_designator_buildingidentifier,
           'locator_designator_buildingidentifierprefix', locator_designator_buildingidentifierprefix,
           'locator_designator_corneraddress1stidentifier', locator_designator_corneraddress1stidentifier,
           'locator_designator_corneraddress2ndidentifier', locator_designator_corneraddress2ndidentifier,
           'locator_designator_entrancedooridentifier', locator_designator_entrancedooridentifier,
           'locator_designator_flooridentifier', locator_designator_flooridentifier,
           'locator_designator_kilometrepoint', locator_designator_kilometrepoint,
           'locator_designator_postaldeliveryidentifier', locator_designator_postaldeliveryidentifier,
           'locator_designator_staircaseidentifier', locator_designator_staircaseidentifier,
           'locator_designator_unitidentifier', locator_designator_unitidentifier,
           'locator_name', locator_name,
           'parcel', parcel,
           'parentaddress', parentaddress,
           'position_specification', position_specification,
           'position_specification_href', position_specification_href,
           'position_method', position_method,
           'position_method_href', position_method_href,
           'position_default', position_default,
           'status', status,
           'status_href', status_href) AS properties,
       ST_Envelope(geom) AS bbox,
       geom
  from addresses.addresses_alternative_encoding;

ALTER TABLE addresses.addresses ADD PRIMARY KEY (fid)
CREATE INDEX addresses_geom_sidx ON addresses.addresses USING GIST (geom);
CREATE INDEX addresses_offsetid_idx ON addresses.addresses(offsetid);
```
