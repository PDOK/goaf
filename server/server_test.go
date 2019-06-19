package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"wfs3_server/codegen"
	"wfs3_server/provider_common"
	gpkg "wfs3_server/provider_gpkg"
)

func TestNewServerWithGeopackageProviderForRoot(t *testing.T) {

	crsMap := make(map[string]string)

	serverEnppoint := "http://testhost:1234"

	gpkgp := gpkg.NewGeopackageProvider(serverEnppoint, "../spec/wfs3.0.yml", "../tst/bgt_wgs84.gpkg", crsMap, "fid", 100, 500)

	server, _ := NewServer(serverEnppoint, "../spec/wfs3.0.yml", 100, 500)
	server, _ = server.SetProviders(gpkgp)

	ts := httptest.NewServer(server.Router())
	defer ts.Close()

	// replace with test endpoint
	gpkgp.ServiceEndpoint = ts.URL

	tests := []struct {
		name  string
		path  string
		want  provider_common.GetLandingPageProvider
		check func(want provider_common.GetLandingPageProvider) error
	}{
		{"root call", "", provider_common.GetLandingPageProvider{}, func(want provider_common.GetLandingPageProvider) error {

			if len(want.Links) != 4 {
				return errors.New("error invalid number of links")
			}

			rels := []string{"self", "service", "conformance", "data"}
			paths := []string{"/", "/api", "/conformance", "/collections"}

			for i, v := range want.Links {
				if v.Rel != rels[i] {
					return errors.New(fmt.Sprintf("Error invalid link rel: %s", v.Rel))
				}

				if v.Href != fmt.Sprintf("%s%s", ts.URL, paths[i]) {
					return errors.New(fmt.Sprintf("Error invalid path rel: %s", v.Href))
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

	serverEnppoint := "http://testhost:1234"

	gpkgp := gpkg.NewGeopackageProvider(serverEnppoint, "../spec/wfs3.0.yml", "../tst/bgt_wgs84.gpkg", crsMap, "fid", 100, 500)

	server, _ := NewServer(serverEnppoint, "../spec/wfs3.0.yml", 100, 500)
	server, _ = server.SetProviders(gpkgp)

	ts := httptest.NewServer(server.Router())
	defer ts.Close()

	// replace with test endpoint
	gpkgp.ServiceEndpoint = ts.URL

	tests := []struct {
		name  string
		path  string
		want  codegen.Content
		check func(want codegen.Content) error
	}{
		{"collection call", "collections", codegen.Content{}, func(want codegen.Content) error {

			if len(want.Collections) != 1 {
				return errors.New(fmt.Sprintf("Error invalid number of collections :%d", len(want.Collections)))
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
