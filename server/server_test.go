package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"oaf-server/codegen"
	"oaf-server/core"
	"oaf-server/geopackage"
	"oaf-server/spec"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	isTesting = true
}

func TestNewServerWithGeopackageProviderForRoot(t *testing.T) {
	serverEndpoint := "http://testhost:1234"

	api, _ := spec.GetOpenAPI("../spec/oaf.json")
	config := core.Config{Datasource: core.Datasource{Geopackage: &core.Geopackage{File: "../example/addresses.gpkg", Fid: "fid"}}}

	gpkgp := geopackage.NewGeopackageWithCommonProvider(api, config)

	server, _ := NewServer(serverEndpoint, "../spec/oaf.json", 100, 500)
	server, _ = server.SetProviders(gpkgp)

	ts := httptest.NewServer(server.Router())
	defer ts.Close()

	// replace with test endpoint
	gpkgp.Config.Service.Url = ts.URL

	tests := []struct {
		name  string
		path  string
		want  core.GetLandingPageProvider
		check func(want core.GetLandingPageProvider) error
	}{
		{"root call", "", core.GetLandingPageProvider{}, func(want core.GetLandingPageProvider) error {

			if len(want.Links) != 11 {
				return errors.New("error invalid number of links")
			}

			rps := map[string][]string{
				"self":         {"?f=json"},
				"alternate":    {"?f=html", "?f=jsonld"},
				"service-doc":  {"/api?f=html"},
				"service-desc": {"/api?f=json"},
				"conformance":  {"/conformance?f=json", "/conformance?f=html", "/conformance?f=jsonld"},
				"data":         {"/collections?f=json", "/collections?f=html", "/collections?f=jsonld"},
			}

			found := false
			for _, v := range want.Links {

				hrefs := rps[v.Rel]
				found = false
				for _, href := range hrefs {

					if v.Href == fmt.Sprintf("%s%s", ts.URL, href) {
						found = true
					}
				}
				if !found {
					return fmt.Errorf("Error invalid path rel: %s", v.Href)
				}
			}

			return nil
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, _ := ts.Client().Get(fmt.Sprintf("%s/%s", ts.URL, tt.path))
			if resp.StatusCode != http.StatusOK {
				t.Fail()
			}
			data, _ := ioutil.ReadAll(resp.Body)

			err := json.Unmarshal(data, &tt.want)
			if err != nil {
				t.Fatal(err)
			}
			err = tt.check(tt.want)
			if err != nil {
				t.Fatal(err)
			}

		})
	}

}

func TestNewServerWithGeopackageProviderForCollection(t *testing.T) {
	serverEndpoint := "http://testhost:1234"

	api, _ := spec.GetOpenAPI("../spec/oaf.json")
	config := core.Config{Datasource: core.Datasource{Geopackage: &core.Geopackage{File: "../example/addresses.gpkg", Fid: "fid"}}}

	gpkgp := geopackage.NewGeopackageWithCommonProvider(api, config)

	server, _ := NewServer(serverEndpoint, "../spec/oaf.json", 100, 500)
	server, _ = server.SetProviders(gpkgp)

	ts := httptest.NewServer(server.Router())
	defer ts.Close()

	// replace with test endpoint
	gpkgp.Config.Service.Url = ts.URL

	tests := []struct {
		name  string
		path  string
		want  codegen.Collections
		check func(want codegen.Collections) error
	}{
		{"collection call", "collections", codegen.Collections{}, func(want codegen.Collections) error {
			if len(want.Collections) != 1 {
				return fmt.Errorf("Error invalid number of collections :%d", len(want.Collections))
			}
			return nil
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, _ := ts.Client().Get(fmt.Sprintf("%s/%s", ts.URL, tt.path))
			if resp.StatusCode != http.StatusOK {
				t.Fail()
			}
			data, _ := ioutil.ReadAll(resp.Body)

			err := json.Unmarshal(data, &tt.want)
			if err != nil {
				t.Fatal(err)
			}
			err = tt.check(tt.want)
			if err != nil {
				t.Fatal(err)
			}

		})
	}

}
