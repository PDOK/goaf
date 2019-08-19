package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"wfs3_server/codegen"
	gpkg "wfs3_server/provider_gpkg"
	postgis "wfs3_server/provider_postgis"
	"wfs3_server/server"

	"github.com/rs/cors"
)

func main() {

	// TODO parse Flags into struct and separate into private method

	//var featureTables arrayFlags

	bindHost := flag.String("s", "0.0.0.0", "server internal bind address, default; 8080")
	bindPort := flag.Int("p", 8080, "server internal bind address, default; 8080")

	serviceEndpoint := flag.String("endpoint", "http://localhost:8080", "server endpoint for proxy reasons, default; http://localhost:8080")
	serviceSpecPath := flag.String("spec", "spec/wfs3.0.yml", "swagger openapi spec")
	defaultReturnLimit := flag.Int("limit", 100, "limit, default: 100")
	maxReturnLimit := flag.Int("limitmax", 500, "max limit, default: 1000")

	providerName := flag.String("provider", "gpkg", "postgis or gpkg")

	gpkgFilePath := flag.String("gpkg", "", "geopackage path")
	crsMapFilePath := flag.String("crs", "", "crs file path")
	configFilePath := flag.String("config", "", "configfile path")

	featureIdKey := flag.String("featureId", "", "Default feature identification or else first column definition (fid)")

	flag.Parse()

	// stage 1: create server with spec path and limits
	apiServer, err := server.NewServer(*serviceEndpoint, *serviceSpecPath, uint64(*defaultReturnLimit), uint64(*maxReturnLimit))
	if err != nil {
		log.Fatal("Server initialisation error:", err)
	}

	// stage 2: Create providers
	var providers codegen.Providers

	if *providerName == "gpkg" {
		providers = addGeopackageProviders(*serviceEndpoint, *serviceSpecPath, *crsMapFilePath, *gpkgFilePath, *featureIdKey, uint64(*defaultReturnLimit), uint64(*maxReturnLimit))

	} else if *providerName == "postgis" {
		providers = addPostgisProviders(*serviceEndpoint, *serviceSpecPath, *configFilePath, uint64(*defaultReturnLimit), uint64(*maxReturnLimit))
	}

	// stage 3: Add providers, also initialises them
	apiServer, err = apiServer.SetProviders(providers)
	if err != nil {
		log.Fatal("Server initialisation error:", err)
	}

	// stage 4: Prepare routing
	router := apiServer.Router()
	handler := cors.Default().Handler(router)

	// ServerEndpoint can be different from bind address due to routing externally
	bindAddress := fmt.Sprintf("%v:%v", *bindHost, *bindPort)

	// print config
	configProvider, err := json.Marshal(apiServer.Providers)
	log.Println(redactPassword(string(configProvider)))

	log.Print("|")
	log.Printf("| SERVING ON: %s", apiServer.ServiceEndpoint)

	// stage 5: Start server
	if err := http.ListenAndServe(bindAddress, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func redactPassword(connectionStr string) (redacted string) {
	var re = regexp.MustCompile(`password=(.*) `)
	redacted = re.ReplaceAllString(connectionStr, `password=******* `)
	return
}

func addPostgisProviders(serviceEndpoint, serviceSpecPath, configFilePath string, defaultReturnLimit, maxReturnLimit uint64) *postgis.PostgisProvider {
	return postgis.NewPostgisProvider(serviceEndpoint, serviceSpecPath, configFilePath, defaultReturnLimit, maxReturnLimit)
}

func addGeopackageProviders(serviceEndpoint, serviceSpecPath, crsMapFilePath string, gpkgFilePath string, featureIdKey string, defaultReturnLimit, maxReturnLimit uint64) *gpkg.GeoPackageProvider {
	crsMap := make(map[string]string)
	csrMapFile, err := ioutil.ReadFile(crsMapFilePath)
	if err != nil {
		log.Printf("Could not read crsmap file: %s, using default CRS Map", crsMapFilePath)
	} else {
		err := json.Unmarshal(csrMapFile, &crsMap)
		log.Print(crsMap)
		if err != nil {
			log.Printf("Could not unmarshal crsmap file: %s, using default CRS Map", crsMapFilePath)
		}

	}
	return gpkg.NewGeopackageProvider(serviceEndpoint, serviceSpecPath, gpkgFilePath, crsMap, featureIdKey, defaultReturnLimit, maxReturnLimit)
}
