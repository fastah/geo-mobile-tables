// Copyright 2020, Blackbuck Computing Inc.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/olekukonko/tablewriter"
)

var mccmnctable map[string]map[string]string
var ipwhoistable map[string]map[string]string

func init() {

	filename := "country-mccmnc-provider.json"
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Read error with JSON document: %v\n", err)
	}

	err = json.Unmarshal(raw, &mccmnctable)
	if err != nil {
		if serr, ok := err.(*json.SyntaxError); ok {
			log.Fatalf("JSON format 💩 in %s at offset of %d bytes: %v", filename, serr.Offset, serr.Error())
		} else {
			log.Fatalf("JSON parsing 💩 in %s : %v", filename, err)
		}
	}

	if len(mccmnctable) == 0 {
		log.Fatalf("JSON for mccmnctable has zero rules :(\n")
	}

	filename = "ipwhois-networkname.json"
	raw, err = ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Read error with JSON document: %v\n", err)
	}

	err = json.Unmarshal(raw, &ipwhoistable)
	if err != nil {
		if serr, ok := err.(*json.SyntaxError); ok {
			log.Fatalf("JSON format 💩 in %s at offset of %d bytes: %v", filename, serr.Offset, serr.Error())
		} else {
			log.Fatalf("JSON parsing 💩 in %s : %v", filename, err)
		}
	}

	if len(ipwhoistable) == 0 {
		log.Fatalf("JSON for ipwhoistable has zero rules :(\n")
	}
}

func TestCrossRef1(t *testing.T) {
	table1 := tablewriter.NewWriter(os.Stdout)
	table1.SetHeader([]string{"Country", "Carrier in MCC/MNC dict", "WHOIS rule hits", "Carrier in WHOIS dict"})
	table1.SetAutoMergeCellsByColumnIndex([]int{0})
	table1.SetRowLine(true)
	table1.SetCaption(true, "LEFT JOIN: MCC/MNC dict to WHOIS dict cross-ref ")

	for country, countryRules := range mccmnctable {
		//t.Logf("Country %s\n", country)
		iprulecount := 0
		if cipwr, exists := ipwhoistable[country]; exists {
			for _, carrierL := range countryRules {
				iprulecount++
				count := 0
				var _cr string
				for _, carrierR := range cipwr {
					if strings.EqualFold(carrierL, carrierR) {
						count++
						_cr = carrierR
					}
				}
				if count > 0 {
					table1.Append([]string{country, carrierL, fmt.Sprintf("%d WHOIS rules", count), _cr})
				} else {
					table1.Append([]string{country, carrierL, fmt.Sprintf("%d WHOIS rules", count), "😬"})
				}
			}
		} else {
			for _, carrierL := range countryRules {
				table1.Append([]string{country, carrierL, "<COUNTRY MISSING IN WHOIS>", "💩"})
			}
		}
	}
	table1.Render()
}

func TestCrossRef2(t *testing.T) {
	table2 := tablewriter.NewWriter(os.Stdout)
	table2.SetHeader([]string{"Country", "Carrier in WHOIS dict", "MCC/MNC rule hits", "Carrier in MCC/MNC dict"})
	table2.SetAutoMergeCellsByColumnIndex([]int{0})
	table2.SetRowLine(true)
	table2.SetCaption(true, "RIGHT JOIN: WHOIS dict to MCC/MNC dict cross-ref ")

	for country, countryRules := range ipwhoistable {
		//t.Logf("Country %s\n", country)
		if ctuples, exists := mccmnctable[country]; exists {
			for _, carrierR := range countryRules {
				count := 0
				var _cl string
				for _, carrierL := range ctuples {
					if strings.EqualFold(carrierR, carrierL) {
						count++
						_cl = carrierL
					}
				}
				if count > 0 {
					table2.Append([]string{country, carrierR, fmt.Sprintf("%d MCC/MNC rules", count), _cl})
				} else {
					table2.Append([]string{country, carrierR, fmt.Sprintf("%d MCC/MNC rules", count), "😬"})
				}
			}
		} else {
			for _, carrierR := range countryRules {
				table2.Append([]string{country, carrierR, "<COUNTRY MISSING IN MCC/MNC>", "💩"})
			}
		}
	}
	table2.Render()
}
