// Copyright 2018, Blackbuck Computing Inc

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/fastah/geo-mobile-tables/utils"
)

func TestMCCMNCFormat(t *testing.T) {

	var ruleBook map[string]map[string]string

	filename := "country-mccmnc-provider.json"
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("Read error with JSON document: %v\n", err)
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

	// Country duplicate entries
	duplicates, err := utils.CheckDuplicateKeys(json.NewDecoder(bytes.NewReader(raw)), nil)
	if err != nil {
		t.Fatalf("Error ! %+v\n", err)
	}
	if len(duplicates) > 0 {
		for _, dup := range duplicates {
			t.Logf("Duplicate key: %s\n", dup)
		}
		t.Fatalf("Duplicate keys found in JSON : %d\n", len(duplicates))
	}

	for country, countryRules := range ruleBook {
		t.Logf("Country %s\n", country)
		for mccmncpair, opname := range countryRules {
			t.Logf("\t%s -> %s\n", mccmncpair, opname)
			tuple := strings.Split(mccmncpair, ",")
			if len(tuple) != 2 {
				t.Fatalf("Invalid MCC MNC key pair %s (%s/%s)\n", mccmncpair, country, opname)
			}
		}
	}
}
