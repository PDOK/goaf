FROM scratch 
 FROM golang:1.11-alpine3.8 AS build-env
 
RUN apk update && apk upgrade && \
   apk add --no-cache bash git gcc musl-dev

ENV GO111MODULE=on

WORKDIR /go/src/server

ADD . /go/src/server

#disable crosscompiling
ENV CGO_ENABLED=1

#compile linux only
ENV GOOS=linux

#build the binary with debug information removed
RUN go build -ldflags '-w -s -linkmode external -extldflags -static' -a -installsuffix cgo -o /wfs-server main.go

FROM scratch as service
EXPOSE 8080
WORKDIR /
ENV PATH=/

COPY --from=build-env  /wfs-server /
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /go/src/server/spec /spec

CMD ["/wfs-server","-s","0.0.0.0","-p","8080", "-gpkg","/srv/data/data.gpkg","-endpoint","http://localhost:8080"]