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
	"wfs3_server/provider_common"
	gpkg "wfs3_server/provider_gpkg"
	postgis "wfs3_server/provider_postgis"
	"wfs3_server/server"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/rs/cors"
)

func main() {

	bindHost := flag.String("s", envString("BIND_HOST", "0.0.0.0"), "server internal bind address, default; 0.0.0.0")
	bindPort := flag.Int("p", envInt("BIND_PORT", 8080), "server internal bind address, default; 8080")

	serviceEndpoint := flag.String("endpoint", envString("ENDPOINT", "http://localhost:8080"), "server endpoint for proxy reasons, default; http://localhost:8080")
	serviceSpecPath := flag.String("spec", envString("SERVICE_SPEC_PATH", "spec/wfs1.0.0.json"), "swagger openapi spec")
	defaultReturnLimit := flag.Int("limit", envInt("LIMIT", 100), "limit, default: 100")
	maxReturnLimit := flag.Int("limitmax", envInt("LIMIT_MAX", 500), "max limit, default: 1000")
	providerName := flag.String("provider", envString("PROVIDER", ""), "postgis or gpkg")
	gpkgFilePath := flag.String("gpkg", envString("PATH_GPKG", ""), "geopackage path")
	crsMapFilePath := flag.String("crs", envString("PATH_CRS", ""), "crs file path")
	configFilePath := flag.String("config", envString("PATH_CONFIG", ""), "configfile path")
	connectionStr := flag.String("connection", envString("CONNECTION", ""), "connection string postgis")
	// alternative database configuration
	if *connectionStr == "" && *providerName == "postgis" {
		withDBHost := flag.String("db-host", envString("DB_HOST", "bgt-cloud-db.postgres.database.azure.com"), "database host")
		withDBPort := flag.Int("db-port", envInt("DB_PORT", 5432), "database port number")
		WithDBName := flag.String("db-name", envString("DB_NAME", "pdok"), "database name")
		withDBSSL := flag.String("db-ssl-mode", envString("DB_SSL_MODE", "disable"), "ssl-mode")
		withDBUser := flag.String("db-user-name", envString("DB_USERNAME", "postgres"), "database username")
		withDBPassword := flag.String("db-password", envString("DB_PASSWORD", ""), "database password")

		connectionStrAlt := fmt.Sprintf("host=%s port=%d database=%s sslmode=%s user=%s password=%s",
			*withDBHost, *withDBPort, *WithDBName, *withDBSSL, *withDBUser, *withDBPassword)

		connectionStr = &connectionStrAlt
	}

	featureIdKey := flag.String("featureId", envString("FEATURE_ID", ""), "Default feature identification or else first column definition (fid)")

	flag.Parse()

	// stage 1: create server with spec path and limits
	apiServer, err := server.NewServer(*serviceEndpoint, *serviceSpecPath, uint64(*defaultReturnLimit), uint64(*maxReturnLimit))
	if err != nil {
		log.Fatal("Server initialisation error:", err)
	}

	// stage 2: Create providers based upon provider name
	commonProvider := provider_common.NewCommonProvider(*serviceEndpoint, *serviceSpecPath, uint64(*defaultReturnLimit), uint64(*maxReturnLimit))
	providers := getProvider(apiServer.Openapi, providerName, commonProvider, crsMapFilePath, gpkgFilePath, featureIdKey, configFilePath, connectionStr)

	if providers == nil {
		log.Fatal("Incorrect provider provided valid names are: gpkg, postgis")
	}

	// stage 3: Add providers, also initialises them
	apiServer, err = apiServer.SetProviders(providers)
	if err != nil {
		log.Fatal("Server initialisation error:", err)
	}

	// stage 4: Prepare routing
	router := apiServer.Router()

	// extra routing for healthcheck
	addHealthHandler(router)

	fs := http.FileServer(http.Dir("/swagger-ui"))
	router.Handler(regexp.MustCompile("/swagger-ui"), http.StripPrefix("/swagger-ui/", fs))

	// cors handler
	handler := cors.Default().Handler(router)

	// ServerEndpoint can be different from bind address due to routing externally
	bindAddress := fmt.Sprintf("%v:%v", *bindHost, *bindPort)

	log.Print("|\n")
	log.Printf("| SERVING ON: %s \n", apiServer.ServiceEndpoint)

	// stage 5: Start server
	if err := http.ListenAndServe(bindAddress, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func getProvider(api *openapi3.T, providerName *string, commonProvider provider_common.CommonProvider, crsMapFilePath *string, gpkgFilePath *string, featureIdKey *string, configFilePath *string, connectionStr *string) codegen.Providers {
	if *providerName == "gpkg" {
		return addGeopackageProviders(api, commonProvider, *crsMapFilePath, *gpkgFilePath, *featureIdKey)
	}
	if *providerName == "postgis" {
		return postgis.NewPostgisWithCommonProvider(api, commonProvider, *configFilePath, *connectionStr)
	}
	return nil
}

func addHealthHandler(router *server.RegexpHandler) {
	router.HandleFunc(regexp.MustCompile("/health"), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write([]byte("ok"))
		if err != nil {
			log.Printf("Could not write ok")
		}
	})
}

func addGeopackageProviders(api *openapi3.T, commonProvider provider_common.CommonProvider, crsMapFilePath string, gpkgFilePath string, featureIdKey string) *gpkg.GeoPackageProvider {
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
	return gpkg.NewGeopackageWithCommonProvider(api, commonProvider, gpkgFilePath, crsMap, featureIdKey)
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
