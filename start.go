package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"wfs3_server/codegen"
	gpkg "wfs3_server/provider_gpkg"
	postgis "wfs3_server/provider_postgis"
	"wfs3_server/server"

	"github.com/rs/cors"
)

func main() {

	bindHost := flag.String("s", envString("BIND_HOST", "0.0.0.0"), "server internal bind address, default; 0.0.0.0")
	bindPort := flag.Int("p", envInt("BIND_PORT", 8080), "server internal bind address, default; 8080")

	serviceEndpoint := flag.String("endpoint", envString("ENDPOINT", "http://localhost:8080"), "server endpoint for proxy reasons, default; http://localhost:8080")
	serviceSpecPath := flag.String("spec", envString("SERVICE_SPEC_PATH", "spec/wfs3.0.yml"), "swagger openapi spec")
	defaultReturnLimit := flag.Int("limit", envInt("LIMIT", 100), "limit, default: 100")
	maxReturnLimit := flag.Int("limitmax", envInt("LIMIT_MAX", 500), "max limit, default: 1000")
	providerName := flag.String("provider", envString("PROVIDER", ""), "postgis or gpkg")
	gpkgFilePath := flag.String("gpkg", envString("PATH_GPKG", ""), "geopackage path")
	crsMapFilePath := flag.String("crs", envString("PATH_CRS", ""), "crs file path")
	configFilePath := flag.String("config", envString("PATH_CONFIG", ""), "configfile path")

	featureIdKey := flag.String("featureId", envString("FEATURE_ID", ""), "Default feature identification or else first column definition (fid)")

	flag.Parse()

	// stage 1: create server with spec path and limits
	apiServer, err := server.NewServer(*serviceEndpoint, *serviceSpecPath, uint64(*defaultReturnLimit), uint64(*maxReturnLimit))
	if err != nil {
		log.Fatal("Server initialisation error:", err)
	}

	// stage 2: Create providers based upon provider name
	var providers codegen.Providers
	if *providerName == "" {
		log.Fatal("No provider provided gpkg/postgis")
	}

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

	// print config with redacted password
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

func envString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}

func envInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value != "" {
		i, e := strconv.ParseInt(value, 10, 32)
		if e != nil {
			return defaultValue
		}
		return int(i)
	}

	return defaultValue
}

func envBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value != "" {
		b, e := strconv.ParseBool(value)
		if e != nil {
			return false
		}
		return b
	}

	return defaultValue
}
