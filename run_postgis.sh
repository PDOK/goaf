#!/bin/bash

docker build -t pdok/wfs3.0:latest .

docker run \
-v $(pwd)/config:/config \
-e CONNECTION='host=127.0.0.1 port=5432 database=bgt_v1 user=postgres password=postgres sslmode=disable' \
-e PROVIDER='postgis' \
-e PATH_CONFIG='/config/config_postgis.yaml' \
-e ENDPOINT='http://localhost:8080' \
-p 8080:8080 \
pdok/wfs3.0:latest
