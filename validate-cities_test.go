// Copyright 2018, Blackbuck Computing Inc

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

func TestCitiesJsonFormat(t *testing.T) {

	raw, err := ioutil.ReadFile("cities-bbox-masterlist.json")
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

	// 1. Counting number of countries and cities in each
	distinctCountries := make(map[string]int)
	for _, city := range cities.Features {
		country := city.Properties["country"].(string)
		cityName := city.Properties["city"].(string)
		if country == "" || cityName == "" {
			t.Errorf("Country or city are MISSING %+v\n", *city)
		}
		distinctCountries[country] = distinctCountries[country] + 1
	}
	//t.Logf("Cities in EACH COUNTRY %+v\n", distinctCountries)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Country", "Cities"})
	for k, v := range distinctCountries {
		table.Append([]string{k, strconv.Itoa(v)})
	}
	table.Render() // Send output

	// 2. Parsing and dimension checks
	for _, city := range cities.Features {

		if !(city.Geometry.GeoJSONType() == "Polygon") {
			t.Errorf("Feature type is not Polygon, must be one. %+v\n", *city)
		}
		if city.Geometry.Dimensions() != 2 {
			t.Errorf("Feature geometry must be a rectangle with NW corner and SE corner coordinates %+v\n", *city)
		}
	}

	// 3. Area checks on bounding boxes of city Polygons
	areasMap := make(map[string]float64)
	for _, city := range cities.Features {
		_, area := planar.CentroidArea(city.Geometry.Bound())
		fakeArea := 10000 * area
		cityName := city.Properties["city"].(string)
		//t.Logf("%s -> area = %f\n", cityName, fakeArea)
		if _, exists := areasMap[cityName]; exists {
			t.Errorf("Duplication entry for city %s\n", cityName)
		} else {
			areasMap[cityName] = fakeArea
		}
	}

	type areakv struct {
		Key   string
		Value float64
	}

	var ss []areakv
	for k, v := range areasMap {
		ss = append(ss, areakv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	// Initialize formatters for country -> city+area viz
	countryAreaTables := make(map[string]*tablewriter.Table)
	for k, _ := range distinctCountries {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Country", "Cities"})
		countryAreaTables[k] = table
	}

	for _, areakv := range ss {
		fmt.Printf("%s, %.1f\n", areakv.Key, areakv.Value)
	}

}
