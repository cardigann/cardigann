package indexer

import (
	"io"
	"net/http"
	"time"

	"github.com/cardigann/cardigann/torznab"
)

// ExampleIndexer provides an example of an indexer for development
type ExampleIndexer struct{}

func (i *ExampleIndexer) Info() torznab.Info {
	return torznab.Info{
		ID:          "example",
		Title:       "Example",
		Description: "An example indexer for testing",
		Link:        "https://example.org/",
		Language:    "en-us",
	}
}

func (i *ExampleIndexer) Test() error {
	return nil
}

func (i *ExampleIndexer) Capabilities() torznab.Capabilities {
	return torznab.Capabilities{
		Categories: torznab.CategoryMapping{
			1: torznab.CategoryMovies,
			2: torznab.CategoryTV,
		},
		SearchModes: []torznab.SearchMode{
			{"search", true, []string{"q"}},
			{"tv-search", true, []string{"q", "season", "ep"}},
		},
	}
}

func (i *ExampleIndexer) Search(query torznab.Query) ([]torznab.ResultItem, error) {
	return []torznab.ResultItem{
		{
			Site:            "example",
			Title:           "Llama llama",
			GUID:            "https://archive.org/download/mma_llama_309960/mma_llama_309960_archive.torrent",
			Link:            "https://archive.org/download/mma_llama_309960/mma_llama_309960_archive.torrent",
			Size:            1600,
			Seeders:         10,
			Peers:           10,
			Category:        torznab.CategoryMovies.ID,
			MinimumRatio:    1,
			MinimumSeedTime: time.Hour * 48,
		},
	}, nil
}

func (i *ExampleIndexer) Download(u string) (io.ReadCloser, http.Header, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, http.Header{}, err
	}

	return resp.Body, resp.Header, nil
}
