package main

import (
	"testing"
)

// camelTests maps db style field_names to camelCase fieldNames
var camelTests = map[string]string{
	"snake_case":  "snakeCase",
	"te":          "te",
	"id":          "ID",
	"manager_id":  "managerID",
	"id_prefix":   "IDPrefix",
	"suffix_html": "suffixHTML",
}

// TestToCamel tests the ToCamel function
func TestToCamel(t *testing.T) {
	for k, v := range camelTests {
		if ToCamel(k) != v {
			t.Fatalf("Failed to convert to camel:%s to:%s result:'%s'", k, v, ToCamel(k))
		}
	}
}

// pluralTests maps singulars to plurals (not reversible)
var pluralTests = map[string]string{
	"man":    "men",
	"person": "people",
	"page":   "pages",
}

// TestPlurals tests the ToPlural function
func TestPlurals(t *testing.T) {
	for k, v := range pluralTests {
		if ToPlural(k) != v {
			t.Fatalf("Failed to convert plural:%s to:%s result:'%s'", k, v, ToPlural(k))
		}
	}
}
