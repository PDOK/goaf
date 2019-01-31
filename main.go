package main

import (
	"flag"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
	gpkg "wfs3_server/provider_gpkg"
	"wfs3_server/server"
)

func main() {

	bindHost := flag.String("s", "0.0.0.0", "server internal bind address, default; 8080")
	bindPort := flag.Int("p", 8080, "server internal bind address, default; 8080")
	serverEndpoint := flag.String("endpoint", "http://localhost:8080", "server endpoint for proxy reasons, default; http://localhost:8080")

	serviceSpecPath := flag.String("spec", "/spec/wfs3.0.yml", "swagger openapi spec")
	filePath := flag.String("gpkg", "/2019_gemeentegrenzen_kustlijn.gpkg", "geopackage path")
	defaultLimit := flag.Int("limit", 20, "limit, default: 20")
	maxLimit := flag.Int("limitmax", 100, "max limit, default: 100")

	flag.Parse()

	apiServer, err := server.NewServerWithGeopackageProvider(&gpkg.GeoPackageProvider{
		ServerEndpoint:  *serverEndpoint,
		ServiceSpecPath: *serviceSpecPath,
		FilePath:        *filePath,
		DefaultLimit:    uint64(*defaultLimit),
		MaxLimit:        uint64(*maxLimit),
	})

	if err != nil {
		log.Fatal("Server initialisation error:", err)
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
