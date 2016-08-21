package indexer

import (
	"reflect"
	"testing"

	"github.com/cardigann/cardigann/torznab"
)

const exampleDefinition1 = `
---
  site: testsite
  name: Test Site
  links:
    - https://www.example.org

  caps:
    categories:
      2:  Movies/BluRay
      6:  Audio
      7:  Movies
      8:  Audio/Video
      10: TV

    modes:
      search: q
      tv-search: [q, season, ep]

  login:
    path: /login.php
    form: form
    inputs:
      username: $username
      password: $password
    error:
      path: /login.php
      message:
        selector: table.detail .text

  search:
    path: torrents.php
    inputs:
      search: $keywords
      cat: 0
    rows:
      selector: table[width='750'] > tbody tr
`

func TestIndexerParser(t *testing.T) {
	def, err := ParseDefinition([]byte(exampleDefinition1))
	if err != nil {
		t.Fatal(err)
	}

	// check defaults
	if def.Language != "en-us" {
		t.Fatalf("Expected language to get the default, got %q", def.Language)
	}

	ok, supported := torznab.Capabilities(def.Capabilities).HasSearchMode("tv-search")
	if !ok {
		t.Fatal("Capabilities should support tv-search")
	}

	if !reflect.DeepEqual(supported, []string{"q", "season", "ep"}) {
		t.Fatalf("Supported parameters for tv-search were parsed incorrectly as %v", supported)
	}

	cat, ok := torznab.Capabilities(def.Capabilities).Categories[6]
	if !ok {
		t.Fatalf("Failed to find a mapping for category 6")
	}

	if cat != torznab.CategoryAudio {
		t.Fatalf("Failed to find a mapping for category 6 to torznab.CategoryAudio")
	}
}
