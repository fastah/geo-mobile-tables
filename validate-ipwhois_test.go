// Copyright 2018, Blackbuck Computing Inc
// portions from StackOverflow with attribution (see code comments)

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func TestIPWhoisFormat(t *testing.T) {

	var ruleBook map[string]map[string]string

	filename := "ipwhois-networkname.json"
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
	duplicates, err := checkDuplicateKeys(json.NewDecoder(bytes.NewReader(raw)), nil)
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
		t.Logf("Observing country %s -> %v\n", country, countryRules)

		if len(country) > 2 {
			t.Fatalf("Country code is greater than 2 characters : %s ( len = %d)\n", country, len(country))
		}
	}
}

// Written by @ThunderCat here: https://stackoverflow.com/questions/50107569/detect-duplicate-in-json-string-golang
// Modified to roll-up all duplicate records up to caller via an additional return value
func checkDuplicateKeys(d *json.Decoder, path []string) (duplicates []string, err error) {
	duplicates = make([]string, 0)
	keys := make(map[string]bool)
	// Get next token from JSON
	t, err := d.Token()
	if err != nil {
		return duplicates, err
	}

	delim, ok := t.(json.Delim)

	// There's nothing to do for simple values (strings, numbers, bool, nil)
	if !ok {
		return duplicates, nil
	}

	switch delim {
	case '{':
		for d.More() {
			// Get field key
			t, err := d.Token()
			if err != nil {
				return duplicates, err
			}
			key := t.(string)

			// Check for duplicates
			if keys[key] {
				//fmt.Printf("Duplicate %s\n", strings.Join(append(path, key), "/"))
				duplicates = append(duplicates, strings.Join(append(path, key), "/"))
			}
			keys[key] = true

			// Check value
			if duplist, err := checkDuplicateKeys(d, append(path, key)); err != nil {
				return duplist, err
			} else {
				// Cumulate duplicate list up the recursion stack
				for _, r := range duplist {
					duplicates = append(duplicates, r)
				}
			}

		}
		// Consume trailing }
		if _, err := d.Token(); err != nil {
			return duplicates, err
		}

	case '[':
		i := 0
		for d.More() {
			if duplist, err := checkDuplicateKeys(d, append(path, strconv.Itoa(i))); err != nil {
				return duplist, err
			} else {
				// Cumulate duplicate list up the recursion stack
				for _, r := range duplist {
					duplicates = append(duplicates, r)
				}
			}

			i++
		}
		// Consume trailing ]
		if _, err := d.Token(); err != nil {
			return duplicates, err
		}

	}
	return duplicates, nil
}
