package utils

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Written by @ThunderCat here: https://stackoverflow.com/questions/50107569/detect-duplicate-in-json-string-golang
// Modified to roll-up all duplicate records up to caller via an additional return value
func CheckDuplicateKeys(d *json.Decoder, path []string) (duplicates []string, err error) {
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
			if duplist, err := CheckDuplicateKeys(d, append(path, key)); err != nil {
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
			if duplist, err := CheckDuplicateKeys(d, append(path, strconv.Itoa(i))); err != nil {
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
