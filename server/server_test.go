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
	gpkg "wfs3_server/provider_gpkg"
)

func TestNewServerWithGeopackageProviderForRoot(t *testing.T) {

	gpkgp := &gpkg.GeoPackageProvider{
		ServerEndpoint:  "http://testhost:1234",
		ServiceSpecPath: "../spec/wfs3.0.yml",
		FilePath:        "../2019_gemeentegrenzen_kustlijn.gpkg",
		//FilePath:        "/media/sf_development-virtual/brtachtergrondkaart.gpkg",l
		//FilePath:     "/media/sf_development-virtual/natura2000.gpkg",
		DefaultLimit: 100,
		MaxLimit:     1000,
	}

	server, _ := NewServerWithGeopackageProvider(gpkgp)
	ts := httptest.NewServer(server.Router())
	defer ts.Close()
	gpkgp.ServerEndpoint = ts.URL

	tests := []struct {
		name  string
		path  string
		want  []codegen.Link
		check func(want []codegen.Link) error
	}{
		{"root call", "", []codegen.Link{}, func(want []codegen.Link) error {

			if len(want) != 4 {
				return errors.New("error invalid number of links")
			}

			rels := []string{"self", "service", "conformance", "data"}
			paths := []string{"/", "/api", "/conformance", "/collections"}

			for i, v := range want {
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

	gpkgp := &gpkg.GeoPackageProvider{
		ServerEndpoint:  "http://testhost:1234",
		ServiceSpecPath: "../spec/wfs3.0.yml",
		FilePath:        "../2019_gemeentegrenzen_kustlijn.gpkg",
		//FilePath:        "/media/sf_development-virtual/brtachtergrondkaart.gpkg",l
		//FilePath:     "/media/sf_development-virtual/natura2000.gpkg",
		DefaultLimit: 20,
		MaxLimit:     100,
	}

	server, _ := NewServerWithGeopackageProvider(gpkgp)
	ts := httptest.NewServer(server.Router())
	defer ts.Close()
	gpkgp.ServerEndpoint = ts.URL

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
