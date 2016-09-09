package indexer

import (
	"reflect"
	"testing"
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

	ok, supported := def.Capabilities.ToTorznab().HasSearchMode("tv-search")
	if !ok {
		t.Fatal("Capabilities should support tv-search")
	}

	if !reflect.DeepEqual(supported, []string{"q", "season", "ep"}) {
		t.Fatalf("Supported parameters for tv-search were parsed incorrectly as %v", supported)
	}

	if l := len(def.Capabilities.ToTorznab().Categories); l != 5 {
		t.Fatalf("Expected 6 categories, got %d", l)
	}
}

const exampleDefinitionWithStringCats = `
---
  site: testsite
  name: Test Site
  links:
    - https://www.example.org

  caps:
    categories:
      abc:  Movies/BluRay
      qyz:  Audio

    modes:
      search: q
      tv-search: [q, season, ep]
`

func TestIndexerParserWithStringLocalCats(t *testing.T) {
	_, err := ParseDefinition([]byte(exampleDefinitionWithStringCats))
	if err != nil {
		t.Fatal(err)
	}
}
