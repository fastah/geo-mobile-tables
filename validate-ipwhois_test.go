// Copyright 2018, Blackbuck Computing Inc

package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestCitiesJsonFormat(t *testing.T) {

	var ruleBook map[string]map[string]string

	filename := "ipwhois-networkname.json"
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("HTTP body read error with GeoJSON document: %v\n", err)
	}

	err = json.Unmarshal(raw, &ruleBook)
	if err != nil {
		if serr, ok := err.(*json.SyntaxError); ok {
			t.Fatalf("JSON format 💩 in %s at offset of %d bytes: %v", filename, serr.Offset, serr.Error())
		} else {
			t.Fatalf("JSON parsing 💩 in %s : %v", filename, err)
		}
	}

	if len(ruleBook) == 0 {
		t.Fatalf("JSON has zero rules :(\n")
	}

	for country, countryRules := range ruleBook {
		t.Fatalf("Observing country %s -> %v\n", country, countryRules)
	}
}
