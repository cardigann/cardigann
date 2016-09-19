package indexer

import (
	"errors"
	"io"
	"net/http"

	"github.com/cardigann/cardigann/torznab"
	"golang.org/x/sync/errgroup"
)

type Aggregate []torznab.Indexer

func (ag Aggregate) Search(query torznab.Query) ([]torznab.ResultItem, error) {
	g := errgroup.Group{}
	allResults := make([][]torznab.ResultItem, len(ag))
	maxLength := 0

	// fetch all results
	for idx, indexer := range ag {
		indexerID := indexer.Info().ID
		idx, indexer := idx, indexer
		g.Go(func() error {
			result, err := indexer.Search(query)
			if err != nil {
				log.Warnf("Indexer %q failed: %s", indexerID, err)
				return nil
			}
			allResults[idx] = result
			if l := len(result); l > maxLength {
				maxLength = l
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		log.Warn(err)
		return nil, err
	}

	results := []torznab.ResultItem{}

	// interleave search results to preserve ordering
	for i := 0; i <= maxLength; i++ {
		for _, r := range allResults {
			if len(r) > i {
				results = append(results, r[i])
			}
		}
	}

	if query.Limit > 0 {
		results = results[:query.Limit]
	}

	return results, nil
}

func (ag Aggregate) Info() torznab.Info {
	return torznab.Info{
		ID:       "aggregate",
		Title:    "Aggregated Indexer",
		Language: "en-US",
		Link:     "",
	}
}

func (ag Aggregate) Capabilities() torznab.Capabilities {
	return torznab.Capabilities{
		SearchModes: []torznab.SearchMode{
			{Key: "tv-search", Available: true, SupportedParams: []string{"q", "season", "ep"}},
			{Key: "search", Available: true, SupportedParams: []string{"q"}},
		},
	}
}

func (ag Aggregate) Download(u string) (io.ReadCloser, http.Header, error) {
	return nil, http.Header{}, errors.New("Not implemented")
}
