package indexer

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/cardigann/cardigann/logger"
	"github.com/cardigann/cardigann/torznab"
)

var (
	log = logger.Logger
)

type TesterOpts struct {
	Download bool
}

type Tester struct {
	Runner *Runner
	Opts   TesterOpts
}

func (t *Tester) Test() error {
	for _, mode := range t.Runner.Capabilities().SearchModes {
		query := torznab.Query{
			Type:  mode.Key,
			Limit: 5,
		}

		switch mode.Key {
		case "tv-search":
			query.Categories = []int{
				torznab.CategoryTV_HD.ID,
				torznab.CategoryTV_SD.ID,
			}
		}

		log.Infof("Testing search mode %q", mode.Key)
		results, err := t.Runner.Search(query)
		if err != nil {
			return err
		}
		if len(results) == 0 {
			return fmt.Errorf("Search returned no results, check logs for details")
		}

		for idx, result := range results {
			if result.Title == "" {
				return fmt.Errorf("Result row %d has empty title", idx+1)
			}
			if result.Size == 0 {
				return fmt.Errorf("Result row %d has zero size", idx+1)
			}
			if result.Link == "" {
				return fmt.Errorf("Result row %d has blank link", idx+1)
			}
			if result.Site == "" {
				return fmt.Errorf("Result row %d has blank site", idx+1)
			}
			if result.Link == "" {
				return fmt.Errorf("Result row %d has empty link", idx+1)
			}
		}

		if t.Opts.Download {
			log.WithField("url", results[0].Link).Infof("Testing downloading torrent")
			rc, _, err := t.Runner.Download(results[0].Link)
			if err != nil {
				return err
			}
			defer rc.Close()

			n, err := io.Copy(ioutil.Discard, rc)
			if err != nil {
				return err
			}

			log.Infof("Downloaded %d bytes from linked torrent", n)
		}
	}

	return nil
}
