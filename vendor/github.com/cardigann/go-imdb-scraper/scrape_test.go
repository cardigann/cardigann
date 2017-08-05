package imdbscraper

import "testing"

func TestScrapingMovie(t *testing.T) {
	m, err := FindByID("tt0087182")
	if err != nil {
		t.Fatal(err)
	}

	if m.Title != "Dune" {
		t.Fatalf("Expected Dune, got %q", m.Title)
	}

	if m.Year != "1984" {
		t.Fatalf("Expected 1984, got %q", m.Year)
	}
}
