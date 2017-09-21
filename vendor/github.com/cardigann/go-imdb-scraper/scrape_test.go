package imdbscraper

import (
	"fmt"
	"testing"
)

func TestScrapingMovies(t *testing.T) {
	testCases := []struct {
		id    string
		title string
		year  string
	}{
		{"tt0087182", "Dune", "1984"},
		{"tt1800302", "Gold", "2016"},
		{"tt0451279", "Wonder Woman", "2017"},
		{"tt0071562", "The Godfather: Part II", "1974"},
		{"tt0045152", "Singin' in the Rain", "1952"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s => %s (%s)", tc.id, tc.title, tc.year), func(t *testing.T) {
			m, err := FindByID(tc.id)
			if err != nil {
				t.Fatal(err)
			}

			if m.Title != tc.title {
				t.Fatalf("Expected %q, got %q", tc.title, m.Title)
			}

			if m.Year != tc.year {
				t.Fatalf("Expected %q, got %q", tc.year, m.Year)
			}
		})
	}
}
