package server

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"log"
	"net/http"
	"wfs3_server/codegen"
	gpkg "wfs3_server/provider_gpkg"
	"wfs3_server/spec"
)

type Server struct {
	Providers codegen.Providers
	swagger   *openapi3.Swagger
}

func NewServerWithGeopackageProvider(providers *gpkg.GeoPackageProvider) (*Server, error) {
	swagger, err := spec.GetSwagger(providers.ServiceSpecPath)

	if err != nil {
		log.Fatal("Specification initialisation error:", err)
		return nil, err
	}

	err = providers.Init()

	if err != nil {
		log.Fatal("Provider initialisation error:", err)
		return nil, err
	}

	return &Server{Providers: providers, swagger: swagger}, nil
}

func (server *Server) HandleForProvider(providerFunc func(r *http.Request) (codegen.Provider, error)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		format, ok := r.URL.Query()["f"]
		if ok && len(format) > 0 {
			r.Header.Add("Layers-Type", format[0])
		}

		provider, err := providerFunc(r)

		if err != nil {
			jsonError(w, "PROVIDER CREATION", err.Error(), http.StatusInternalServerError)
			return
		}

		if provider == nil {
			http.NotFound(w, r)
			return
		}

		ct := r.Header.Get("Layers-Type")

		if ct == "" {
			ct = codegen.JSONContentType
		}

		result, err := provider.Provide()

		if err != nil {
			jsonError(w, "PROVIDER", err.Error(), http.StatusInternalServerError)
			return
		}

		var encodedContent []byte

		if ct == codegen.JSONContentType {
			encodedContent, err = provider.MarshalJSON(result)
			if err != nil {
				jsonError(w, "JSON MARSHALLER", err.Error(), http.StatusInternalServerError)
				return
			}

		} else if ct == codegen.HTMLContentType {
			encodedContent, err = provider.MarshalHTML(result)
			if err != nil {
				jsonError(w, "HTML MARSHALLER", err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			jsonError(w, "Invalid Layers Type", "Layers-Type: ''"+ct+"'' not supported.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Layers-Type", ct)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(encodedContent)
	}
}

func jsonError(w http.ResponseWriter, code string, msg string, status int) {
	w.WriteHeader(status)

	result, err := json.Marshal(&codegen.Exception{
		Code:        code,
		Description: msg,
	})

	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("problem marshaling error: %v", msg)))
	} else {
		_, _ = w.Write(result)
	}
}
