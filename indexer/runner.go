package indexer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cardigann/cardigann/config"
	"github.com/cardigann/cardigann/torznab"
	"github.com/dustin/go-humanize"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
)

var (
	_ torznab.Indexer = &Runner{}
)

type Runner struct {
	Definition *IndexerDefinition
	Browser    browser.Browsable
	Config     config.Config
	caps       torznab.Capabilities
}

func NewRunner(def *IndexerDefinition, conf config.Config) *Runner {
	bow := surf.NewBrowser()
	bow.SetUserAgent(agent.Chrome())
	bow.SetAttribute(browser.SendReferer, false)
	bow.SetAttribute(browser.MetaRefreshHandling, false)

	return &Runner{
		Definition: def,
		Browser:    bow,
		Config:     conf,
	}
}

func (r *Runner) ResolveVariable(name string, resolver func(string) (string, error)) (string, error) {
	if name[0] == '$' {
		return resolver(strings.TrimPrefix(name, "$"))
	}
	return name, nil
}

func (r *Runner) ResolvePath(urlPath string) (string, error) {
	var urlStr string

	if configUrl, ok, _ := r.Config.Get(r.Definition.Site, "url"); ok {
		urlStr = configUrl
	} else {
		urlStr = r.Definition.Links[0]
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	u.Path = urlPath
	return u.String(), nil
}

func (r *Runner) checkLoginError() error {
	if !r.Definition.Login.Error.MatchPage(r.Browser) {
		return nil
	}

	msg, err := r.Definition.Login.Error.Message.Text(r.Browser.Dom())
	if err != nil {
		return err
	}

	return errors.New(strings.TrimSpace(msg))
}

func (r *Runner) Login() error {
	loginUrl, err := r.ResolvePath(r.Definition.Login.Path)
	if err != nil {
		return err
	}

	log.Printf("[%s] Attempting to login to %s",
		r.Definition.Site, loginUrl)

	err = r.Browser.Open(loginUrl)
	if err != nil {
		return err
	}

	log.Printf("[%s] Status code is %d, landed on page %s",
		r.Definition.Site, r.Browser.StatusCode(), r.Browser.Url())

	fm, err := r.Browser.Form(r.Definition.Login.FormSelector)
	if err != nil {
		return err
	}

	for name, val := range r.Definition.Login.Inputs {
		log.Printf("[%s] Filling input %q in form %q with %q",
			r.Definition.Site, name, r.Definition.Login.FormSelector, val)

		resolved, err := r.ResolveVariable(val, func(name string) (string, error) {
			s, _, err := r.Config.Get(r.Definition.Site, strings.TrimPrefix(name, "$"))
			return s, err
		})
		if err != nil {
			return err
		}

		if err = fm.Input(name, resolved); err != nil {
			return err
		}
	}

	log.Printf("[%s] Submitting login form", r.Definition.Site)

	if err = fm.Submit(); err != nil {
		log.Printf("[%s] Login failed with %q", r.Definition.Site, err.Error())
		return err
	}

	log.Printf("[%s] Status code is %d, landed on page %s",
		r.Definition.Site, r.Browser.StatusCode(), r.Browser.Url())

	if err = r.checkLoginError(); err != nil {
		log.Printf("[%s] Failed to login with %q", r.Definition.Site, err.Error())
		return err
	}

	log.Printf("[%s] Successfully logged in", r.Definition.Site)
	return nil
}

func (r *Runner) Info() torznab.Info {
	return torznab.Info{
		ID:       r.Definition.Site,
		Title:    r.Definition.Name,
		Language: r.Definition.Language,
	}
}

func (r *Runner) Test() error {
	for _, mode := range r.Capabilities().SearchModes {
		log.Printf("[%s] Testing search mode %s", r.Definition.Site, mode.Key)
		results, err := r.Search(torznab.Query{"t": mode.Key})
		if err != nil {
			return err
		}
		if len(results) == 0 {
			return torznab.ErrNoSuchItem
		}
	}

	return nil
}

func (r *Runner) Capabilities() torznab.Capabilities {
	return torznab.Capabilities(r.Definition.Capabilities)
}

func (r *Runner) Search(query torznab.Query) ([]torznab.ResultItem, error) {
	searchUrl, err := r.ResolvePath(r.Definition.Search.Path)
	if err != nil {
		return nil, err
	}

	log.Printf("[%s] Opening %s", r.Definition.Site, searchUrl)

	err = r.Browser.Open(searchUrl)
	if err != nil {
		return nil, err
	}

	log.Printf("[%s] Status code is %d, landed on page %s",
		r.Definition.Site, r.Browser.StatusCode(), r.Browser.Url())

	// log.Println(r.Browser.ResponseHeaders())
	// log.Println(r.Browser.Body())

	vals := url.Values{}

	for name, val := range r.Definition.Search.Inputs {
		resolved, err := r.ResolveVariable(val, func(name string) (string, error) {
			switch name {
			case "keywords":
				return query.Keywords(), nil
			}
			return "", errors.New("Undefined variable " + name)
		})
		if err != nil {
			return nil, err
		}
		vals.Add(name, resolved)
	}

	log.Printf("[%s] Query params %s", r.Definition.Site, vals.Encode())

	err = r.Browser.OpenForm(searchUrl, vals)
	if err != nil {
		return nil, err
	}

	log.Printf("[%s] Status code is %d, landed on page %s",
		r.Definition.Site, r.Browser.StatusCode(), r.Browser.Url())

	items := []torznab.ResultItem{}
	timer := time.Now()
	rows := r.Browser.Find(r.Definition.Search.Rows.Selector)

	log.Printf("[%s] Found %d rows matching %q",
		r.Definition.Site, rows.Length(), r.Definition.Search.Rows.Selector)

	for i := 0; i < rows.Length(); i++ {
		row := map[string]string{}

		for field, block := range r.Definition.Search.Fields {
			log.Printf("[%s] Processing field %q of row %d (selector %q)",
				r.Definition.Site, field, i+1, block.Selector)

			val, err := block.Text(rows.Eq(i))
			if err != nil {
				return nil, err
			}

			row[field] = val
		}

		item := torznab.ResultItem{
			Site:            r.Definition.Site,
			MinimumRatio:    1,
			MinimumSeedTime: time.Hour * 48,
		}

		log.Printf("[%s] Row %d %+v", r.Definition.Site, i+1, row)

		for key, val := range row {
			switch key {
			case "download":
				u, err := r.Browser.ResolveStringUrl(val)
				if err != nil {
					return nil, fmt.Errorf("Search result row #%d has malformed url in %s", i+1, key)
				}
				item.Link = u
			case "details":
				u, err := r.Browser.ResolveStringUrl(val)
				if err != nil {
					return nil, fmt.Errorf("Search result row #%d has malformed url in %s", i+1, key)
				}
				item.GUID = u
			case "comments":
				u, err := r.Browser.ResolveStringUrl(val)
				if err != nil {
					return nil, fmt.Errorf("Search result row #%d has malformed url in %s", i+1, key)
				}
				item.Comments = u
			case "title":
				item.Title = val
			case "description":
				item.Description = val
			case "category":
				catId, err := strconv.Atoi(val)
				if err != nil {
					return nil, fmt.Errorf("Search result row #%d has malformed category id in %s", i+1, key)
				}
				mappedCat, ok := torznab.Capabilities(r.Definition.Capabilities).Categories[catId]
				if !ok {
					return nil, fmt.Errorf("Search result row #%d has unmappable category id %d in %s", i+1, catId, key)
				}
				item.Category = mappedCat.ID
			case "size":
				bytes, err := humanize.ParseBytes(val)
				if err != nil {
					return nil, fmt.Errorf("Search result row #%d has malformed size: %s", i+1, err.Error())
				}
				item.Size = bytes
			case "leechers":
				leechers, err := strconv.Atoi(val)
				if err != nil {
					return nil, fmt.Errorf("Search result row #%d has malformed leechers value in %s", i+1, key)
				}
				item.Peers += leechers
			case "seeders":
				seeders, err := strconv.Atoi(val)
				if err != nil {
					return nil, fmt.Errorf("Search result row #%d has malformed seeders value in %s", i+1, key)
				}
				item.Seeders = seeders
				item.Peers += seeders
			case "date":
				t, err := time.Parse(time.RFC1123Z, val)
				if err != nil {
					return nil, fmt.Errorf("Search result row #%d has malformed time value in %s", i+1, key)
				}
				item.PublishDate = t
			default:
				return nil, fmt.Errorf("Unknown field %q", key)
			}
		}

		items = append(items, item)
	}

	log.Printf("Found %d results in %s", len(items), time.Now().Sub(timer))
	return items, nil
}

func (r *Runner) Download(u string) (io.ReadCloser, http.Header, error) {
	if err := r.Login(); err != nil {
		return nil, http.Header{}, err
	}

	fullUrl, err := r.Browser.ResolveStringUrl(u)
	if err != nil {
		return nil, http.Header{}, err
	}

	if err := r.Browser.Open(fullUrl); err != nil {
		return nil, http.Header{}, err
	}

	b := &bytes.Buffer{}

	if _, err := r.Browser.Download(b); err != nil {
		return nil, http.Header{}, err
	}

	return ioutil.NopCloser(bytes.NewReader(b.Bytes())), r.Browser.ResponseHeaders(), nil
}
