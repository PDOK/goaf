package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"oaf-server/codegen"
	"oaf-server/provider_common"
	gpkg "oaf-server/provider_gpkg"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestNewServerWithGeopackageProviderForRoot(t *testing.T) {

	crsMap := make(map[string]string)

	serverEndpoint := "http://testhost:1234"

	commonProvider := provider_common.NewCommonProvider(serverEndpoint, "../spec/oaf.json", 100, 500)
	gpkgp := gpkg.NewGeopackageWithCommonProvider(nil, commonProvider, "../example/bgt_wgs84.gpkg", crsMap, "fid")

	server, _ := NewServer(serverEndpoint, "../spec/oaf.json", 100, 500)
	server, _ = server.SetProviders(gpkgp)

	ts := httptest.NewServer(server.Router())
	defer ts.Close()

	// replace with test endpoint
	gpkgp.CommonProvider.ServiceEndpoint = ts.URL

	tests := []struct {
		name  string
		path  string
		want  provider_common.GetLandingPageProvider
		check func(want provider_common.GetLandingPageProvider) error
	}{
		{"root call", "", provider_common.GetLandingPageProvider{}, func(want provider_common.GetLandingPageProvider) error {

			if len(want.Links) != 8 {
				return errors.New("error invalid number of links")
			}

			rels := []string{"self", "alternate", "service", "service", "conformance", "conformance", "data", "data"}
			paths := []string{"?f=json", "?f=html", "/api?f=json", "/api?f=html", "/conformance?f=json", "/conformance?f=html", "/collections?f=json", "/collections?f=html"}

			for i, v := range want.Links {
				if v.Rel != rels[i] {
					return fmt.Errorf("Error invalid link rel: %s", v.Rel)
				}

				if v.Href != fmt.Sprintf("%s%s", ts.URL, paths[i]) {
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

	crsMap := make(map[string]string)

	serverEndpoint := "http://testhost:1234"

	commonProvider := provider_common.NewCommonProvider(serverEndpoint, "../spec/oaf.json", 100, 500)
	gpkgp := gpkg.NewGeopackageWithCommonProvider(nil, commonProvider, "../example/bgt_wgs84.gpkg", crsMap, "fid")

	server, _ := NewServer(serverEndpoint, "../spec/oaf.json", 100, 500)
	server, _ = server.SetProviders(gpkgp)

	ts := httptest.NewServer(server.Router())
	defer ts.Close()

	// replace with test endpoint
	gpkgp.CommonProvider.ServiceEndpoint = ts.URL

	tests := []struct {
		name  string
		path  string
		want  codegen.Collections
		check func(want codegen.Collections) error
	}{
		{"collection call", "collections", codegen.Collections{}, func(want codegen.Collections) error {
			if len(want.Collections) != 27 {
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
