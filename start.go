package main

import (
	"flag"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
	gpkg "wfs3_server/provider_gpkg"
	postgis "wfs3_server/provider_postgis"
	"wfs3_server/server"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {

	var featureTables arrayFlags

	bindHost := flag.String("s", "0.0.0.0", "server internal bind address, default; 8080")
	bindPort := flag.Int("p", 8080, "server internal bind address, default; 8080")
	serverEndpoint := flag.String("endpoint", "http://localhost:8080", "server endpoint for proxy reasons, default; http://localhost:8080")

	serviceSpecPath := flag.String("spec", "spec/wfs3.0.yml", "swagger openapi spec")
	gpkgFilePath := flag.String("gpkg", "", "geopackage path")
	connectionStr := flag.String("postgis", "", "postgis connection str")

	flag.Var(&featureTables, "collection", "postgis feature table, can be repeated multiple times.")
	featureIdKey := flag.String("featureId", "", "Default feature identification or else first column definition (fid)")
	defaultLimit := flag.Int("limit", 20, "limit, default: 20")
	maxLimit := flag.Int("limitmax", 100, "max limit, default: 100")

	flag.Parse()

	var apiServer *server.Server

	if *gpkgFilePath != "" {
		api, err := server.NewServerWithGeopackageProvider(&gpkg.GeoPackageProvider{
			ServerEndpoint:  *serverEndpoint,
			ServiceSpecPath: *serviceSpecPath,
			FilePath:        *gpkgFilePath,
			FeatureIdKey:    *featureIdKey,
			DefaultLimit:    uint64(*defaultLimit),
			MaxLimit:        uint64(*maxLimit),
		})

		if err != nil {
			log.Fatal("Server initialisation error:", err)
		}
		apiServer = api

	} else if *connectionStr != "" {
		api, err := server.NewServerWithPostgisProvider(&postgis.PostgisProvider{
			ServerEndpoint:  *serverEndpoint,
			ServiceSpecPath: *serviceSpecPath,
			ConnectionStr:   *connectionStr,
			FeatureTables:   featureTables,
			DefaultLimit:    uint64(*defaultLimit),
			MaxLimit:        uint64(*maxLimit),
		})

		if err != nil {
			log.Fatal("Server initialisation error:", err)
		}
		apiServer = api

	}

	router := apiServer.Router()
	handler := cors.Default().Handler(router)

	// ServerEndpoint can be different from bindaddress due to routing externally
	bindAddress := fmt.Sprintf("%v:%v", *bindHost, *bindPort)

	log.Print("|")
	log.Printf("| SERVING ON: %s", *serverEndpoint)

	if err := http.ListenAndServe(bindAddress, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}