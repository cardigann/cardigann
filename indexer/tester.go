package indexer

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/cardigann/cardigann/logger"
	"github.com/cardigann/cardigann/torznab"
	"github.com/mgutz/ansi"
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
	Output io.Writer
}

func (t *Tester) printf(format string, args ...interface{}) {
	w := t.Output
	if w == nil {
		w = os.Stdout
	}

	fmt.Fprintf(w, format, args...)
}

func (t *Tester) printfWithResult(format string, args []interface{}, f func() error) error {
	timer := time.Now()
	t.printf(format+" ", args...)

	err := f()
	if err == nil {
		t.printf(" %s %s\n",
			ansi.Color("SUCCESS ✓", "green"),
			ansi.Color("in "+time.Now().Sub(timer).String(), "white"))
	} else {
		t.printf(" %s %s\n",
			ansi.Color("FAILURE ✗", "red"),
			ansi.Color("in "+time.Now().Sub(timer).String(), "white"))
	}

	return err
}

func (t *Tester) testSearchMode(mode torznab.SearchMode) error {
	query := torznab.Query{
		Type:  mode.Key,
		Limit: 3,
	}

	switch mode.Key {
	case "tv-search":
		query.Categories = []int{
			torznab.CategoryTV_HD.ID,
			torznab.CategoryTV_SD.ID,
		}
	}

	results, err := t.Runner.Search(query)
	if err != nil {
		return err
	}

	return t.assertValidResults(results)
}

func (t *Tester) testLogin() error {
	return t.Runner.login()
}

func (t *Tester) assertValidResults(results []torznab.ResultItem) error {
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

		if t.Opts.Download {
			if err := t.assertValidTorrent(result); err != nil {
				return err
			}
		}

	}
	return nil
}

func (t *Tester) assertValidTorrent(result torznab.ResultItem) error {
	u, err := url.Parse(result.Link)
	if err != nil {
		return err
	}

	if u.Scheme == "magnet" {
		return nil
	}

	rc, _, err := t.Runner.Download(result.Link)
	if err != nil {
		return err
	}
	defer rc.Close()

	_, err = io.Copy(ioutil.Discard, rc)
	return err
}

func (t *Tester) Test() error {
	info := t.Runner.Info()
	t.printf("→ Testing indexer %s (%s)\n", info.ID, info.Link)
	var err error

	defer func() {
		if err != nil {
			t.printf("→ Indexer %s %s\n", info.ID, ansi.Color("FAILED", "red"))
		} else {
			t.printf("→ Indexer %s is %s\n", info.ID, ansi.Color("OK", "green"))
		}
	}()

	err = t.printfWithResult("  Testing login with valid credentials", nil, func() error {
		return t.testLogin()
	})

	for _, mode := range t.Runner.Capabilities().SearchModes {
		mode := mode
		err = t.printfWithResult("  Testing search mode %s", []interface{}{mode.Key}, func() error {
			return t.testSearchMode(mode)
		})
		if err != nil {
			return err
		}
	}

	err = t.printfWithResult("  Testing empty results are handled", nil, func() error {
		results, err := t.Runner.Search(torznab.Query{
			Q: "nothingshouldmatchtheseresults",
		})
		if err != nil {
			return err
		}
		if len(results) > 0 {
			return fmt.Errorf("Expected no results, got %d", len(results))
		}
		return nil
	})

	if err != nil {
		return err
	}

	err = t.printfWithResult("  Testing ratio", nil, func() error {
		_, err := t.Runner.Ratio()
		return err
	})

	if err != nil {
		return err
	}

	// log.Infof("Ratio returned %s", ratio)
	return nil
}
