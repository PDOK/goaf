package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"oaf-server/codegen"
	"oaf-server/gpkg"
	"oaf-server/provider"
	"oaf-server/spec"

	"github.com/getkin/kin-openapi/openapi3"
)

type Server struct {
	ServiceEndpoint    string
	ServiceSpecPath    string
	MaxReturnLimit     uint64
	DefaultReturnLimit uint64
	Providers          codegen.Providers
	Openapi            *openapi3.T
	Templates          *template.Template
}

func NewServer(serviceEndpoint, serviceSpecPath string, defaultReturnLimit, maxReturnLimit uint64) (*Server, error) {
	openapi, err := spec.GetOpenAPI(serviceSpecPath)

	if err != nil {
		log.Fatal("Specification initialisation error:", err)
		return nil, err
	}
	// Set endpoint
	openapi.AddServer(&openapi3.Server{URL: serviceEndpoint, Description: "Production server"})

	server := &Server{ServiceEndpoint: serviceEndpoint, ServiceSpecPath: serviceSpecPath, MaxReturnLimit: maxReturnLimit, DefaultReturnLimit: defaultReturnLimit, Openapi: openapi}

	// add templates to server
	server.Templates = template.Must(template.New("templates").Funcs(
		template.FuncMap{
			"isOdd":       func(i int) bool { return i%2 != 0 },
			"hasFeatures": func(i []gpkg.Feature) bool { return len(i) > 0 },
			"upperFirst":  provider.UpperFirst,
			"dict": func(values ...interface{}) (map[string]interface{}, error) {
				if len(values)%2 != 0 {
					return nil, errors.New("invalid dict call")
				}
				dict := make(map[string]interface{}, len(values)/2)
				for i := 0; i < len(values); i += 2 {
					key, ok := values[i].(string)
					if !ok {
						return nil, errors.New("dict keys must be strings")
					}
					dict[key] = values[i+1]
				}
				return dict, nil
			},
			//}).ParseGlob("/templates/*")) // prod
		}).ParseGlob("templates/*")) // IDE

	return server, nil
}

func (s *Server) SetProviders(providers codegen.Providers) (*Server, error) {
	err := providers.Init()

	if err != nil {
		log.Fatal("Provider initialiation error:", err)
		return nil, err
	}
	s.Providers = providers
	return s, nil
}

func (s *Server) HandleForProvider(providerFunc func(r *http.Request) (codegen.Provider, error)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		p, err := providerFunc(r)

		// TODO: improve formatting errors (spec might require specific json schema in response)
		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			switch v := err.(type) {
			default:
				jsonError(w, "PROVIDER CREATION", v.Error(), http.StatusNotFound)
				return
			case *provider.InvalidContentTypeError:
				jsonError(w, "CLIENT ERROR", v.Error(), http.StatusBadRequest)
				return
			case *provider.InvalidFormatError:
				jsonError(w, "CLIENT ERROR", v.Error(), http.StatusBadRequest)
				return
			}
		}

		if p == nil {
			http.NotFound(w, r)
			return
		}

		result, err := p.Provide()

		// todo  error based on content type
		if err != nil {
			jsonError(w, "PROVIDER", err.Error(), http.StatusInternalServerError)
			return
		}

		var encodedContent []byte

		if p.ContentType() == provider.JSONContentType || p.ContentType() == provider.LDJSONContentType || p.ContentType() == provider.GEOJSONContentType {
			encodedContent, err = json.Marshal(result)
			if err != nil {
				jsonError(w, "JSON MARSHALLER", err.Error(), http.StatusInternalServerError)
				return
			}
		} else if p.ContentType() == provider.HTMLContentType {
			providerID := p.String()

			rmap := make(map[string]interface{})
			rmap["result"] = result
			rmap["srsid"] = p.SrsId()

			b := new(bytes.Buffer)
			err = s.Templates.ExecuteTemplate(b, providerID+".html", rmap)
			encodedContent = b.Bytes()

			if err != nil {
				jsonError(w, "HTML MARSHALLER", err.Error(), http.StatusInternalServerError)
				return
			}

		} else {
			jsonError(w, "Invalid Content Type", "Content-Type: ''"+p.ContentType()+"'' not supported.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", p.ContentType())
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
