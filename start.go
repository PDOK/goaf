package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"oaf-server/codegen"
	"oaf-server/gpkg"
	"oaf-server/postgis"
	"oaf-server/provider"
	"oaf-server/server"
	"os"
	"regexp"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v2"

	"github.com/rs/cors"
)

func main() {

	bindHost := flag.String("s", envString("BIND_HOST", "0.0.0.0"), "server internal bind address, default; 0.0.0.0")
	bindPort := flag.Int("p", envInt("BIND_PORT", 8080), "server internal bind address, default; 8080")

	configfilepath := flag.String("c", envString("CONFIG", ""), "configfile path")
	flag.Parse()

	config := &provider.Config{}
	config.ReadConfig(*configfilepath)

	// stage 1: create server with spec path and limits
	apiServer, err := server.NewServer(config.Service.Url, config.Openapi, uint64(config.DefaultFeatureLimit), uint64(config.MaxFeatureLimit))
	if err != nil {
		log.Fatal("Server initialisation error:", err)
	}

	// stage 2: Create providers based upon provider name
	commonProvider := provider.NewCommonProvider(config.Openapi, uint64(config.DefaultFeatureLimit), uint64(config.MaxFeatureLimit))
	// providers := getProvider(apiServer.Openapi, providerName, commonProvider, crsMapFilePath, gpkgFilePath, featureIdKey, configFilePath, connectionStr)

	providers := getProvider(apiServer.Openapi, commonProvider, *config)

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

	fs := http.FileServer(http.Dir("swagger-ui"))
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

func getProvider(api *openapi3.T, commonProvider provider.CommonProvider, config provider.Config) codegen.Providers {
	if config.Datasource.Geopackage != nil {
		return addGeopackageProviders(api, commonProvider, "", config)
	} else if config.Datasource.PostGIS != nil {
		return postgis.NewPostgisWithCommonProvider(api, commonProvider, config)
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

func addGeopackageProviders(api *openapi3.T, commonProvider provider.CommonProvider, crsMapFilePath string, config provider.Config) *gpkg.GeoPackageProvider {
	crsMap := make(map[string]string)
	csrMapFile, err := ioutil.ReadFile(crsMapFilePath)
	if err != nil {
		log.Printf("Could not read crsmap file: %s, using default CRS Map", crsMapFilePath)
	} else {
		err := yaml.Unmarshal(csrMapFile, &crsMap)
		log.Print(crsMap)
		if err != nil {
			log.Printf("Could not unmarshal crsmap file: %s, using default CRS Map", crsMapFilePath)
		}
	}

	if crsMap[`4326`] == `` {
		crsMap[`4326`] = `http://www.opengis.net/def/crs/OGC/1.3/CRS84`
	}

	return gpkg.NewGeopackageWithCommonProvider(api, commonProvider, crsMap, config)
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
