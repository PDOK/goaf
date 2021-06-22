#!/bin/bash

docker build -t pdok/goaf:latest .

docker run \
-v $(pwd)/example:/config \
-e CONNECTION='host=127.0.0.1 port=5432 database=bgt_v1 user=postgres password=postgres sslmode=disable' \
-e PROVIDER='postgis' \
-e PATH_CONFIG='/config/config_postgis.yaml' \
-e ENDPOINT='http://localhost:8080' \
-p 8080:8080 \
pdok/goaf:latest
