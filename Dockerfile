FROM golang:1.16-alpine3.13 AS build-env
 
RUN apk update && apk upgrade && \
   apk add --no-cache bash git gcc musl-dev

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org

WORKDIR /go/src/server
ADD . /go/src/server

RUN go mod download

# crosscompiling
ENV CGO_ENABLED=1

# compile linux only
ENV GOOS=linux

# build the binary with debug information removed
RUN go build -ldflags '-w -s -linkmode external -extldflags -static' -a -installsuffix cgo -o /oaf-server start.go

FROM scratch as service
EXPOSE 8080
WORKDIR /
ENV PATH=/

COPY --from=build-env /go/src/server/spec/oaf.json /spec/oaf.json
COPY --from=build-env /oaf-server /
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /go/src/server/templates /templates
COPY --from=build-env /go/src/server/swagger-ui /swagger-ui

CMD ["/oaf-server"]