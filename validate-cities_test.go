// Copyright 2018, Blackbuck Computing Inc

package main

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/paulmach/orb/geojson"
)

func TestCitiesJsonFormat(t *testing.T) {

	inputUrl := "https://raw.githubusercontent.com/fastah/geo-mobile-tables/master/cities-bbox-masterlist.json"
	// fetch and parse master geojson
	resp, err := http.Get(inputUrl)
	if err != nil {
		t.Errorf("Problem fetching GeoJSON file %v (url = %s)\n", err, inputUrl)
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("HTTP body read error with GeoJSON document: %v\n", err)
	}

	cities, err := geojson.UnmarshalFeatureCollection(raw)
	if err != nil {
		t.Errorf("GeoJSON document doesn't have a valid syntax: %v\n", err)
	}

	if len(cities.Features) == 0 {
		t.Errorf("Features array inside GeoJSON should not be zero :(\n")
	}

	for _, city := range cities.Features {

		if !(city.Geometry.GeoJSONType() == "Polygon") {
			t.Errorf("Feature type is not Polygon, must be one. %+v\n", *city)
		}
		if city.Geometry.Dimensions() != 2 {
			t.Errorf("Feature geometry must be a rectangle with NW corner and SE corner coordinates %+v\n", *city)
		}
	}
}
