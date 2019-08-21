package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"wfs3_server/codegen"
	"wfs3_server/spec"

	"github.com/getkin/kin-openapi/openapi3"
)

type Server struct {
	ServiceEndpoint    string
	ServiceSpecPath    string
	MaxReturnLimit     uint64
	DefaultReturnLimit uint64
	Providers          codegen.Providers
	swagger            *openapi3.Swagger
	Templates          *template.Template
}

func NewServer(serviceEndpoint, serviceSpecPath string, defaultReturnlimit, maxReturnLimit uint64) (*Server, error) {
	swagger, err := spec.GetSwagger(serviceSpecPath)

	if err != nil {
		log.Fatal("Specification initialisation error:", err)
		return nil, err
	}

	server := &Server{ServiceEndpoint: serviceEndpoint, ServiceSpecPath: serviceSpecPath, MaxReturnLimit: maxReturnLimit, DefaultReturnLimit: defaultReturnlimit, swagger: swagger}

	// add templates to server
	server.Templates = template.Must(template.New("templates").Funcs(
		template.FuncMap{}).ParseGlob("templates/*"))

	return server, nil
}

func (s *Server) SetProviders(providers codegen.Providers) (*Server, error) {
	err := providers.Init()

	if err != nil {
		log.Fatal("Provider initialisation error:", err)
		return nil, err
	}

	s.Providers = providers

	return s, nil
}

func (s *Server) HandleForProvider(providerFunc func(r *http.Request) (codegen.Provider, error)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		format, ok := r.URL.Query()["f"]
		if ok && len(format) > 0 {
			r.Header.Add("Content-Type", format[0])
		}

		provider, err := providerFunc(r)

		if err != nil {
			jsonError(w, "PROVIDER CREATION", err.Error(), http.StatusNotFound)
			return
		}

		if provider == nil {
			http.NotFound(w, r)
			return
		}

		ct := r.Header.Get("Content-Type")

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
			encodedContent, err = json.Marshal(result)
			if err != nil {
				jsonError(w, "JSON MARSHALLER", err.Error(), http.StatusInternalServerError)
				return
			}

		} else if ct == codegen.HTMLContentType {
			encodedContent, err = json.Marshal(result)
			if err != nil {
				jsonError(w, "HTML MARSHALLER", err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			jsonError(w, "Invalid Content Type", "Content-Type: ''"+ct+"'' not supported.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", ct)
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
