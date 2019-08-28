#!/bin/bash

docker build -t pdok/wfs3.0:latest .

docker run \
-v $(pwd)/tst:/config \
-e PROVIDER='gpkg' \
-e PATH_GPKG='/config/bgt_wgs84.gpkg' \
-e ENDPOINT='http://localhost:8080' \
-p 8080:8080 \
pdok/wfs3.0:latest
