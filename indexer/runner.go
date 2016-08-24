package indexer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

	log "github.com/Sirupsen/logrus"
)

var (
	_ torznab.Indexer = &Runner{}
)

type Runner struct {
	Definition *IndexerDefinition
	Browser    browser.Browsable
	Config     config.Config
	caps       torznab.Capabilities
	logger     log.FieldLogger
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
		logger: log.New().WithFields(log.Fields{
			"indexer": def.Site,
		}),
	}
}

func (r *Runner) resolveVariable(name string, resolver func(string) (string, error)) (string, error) {
	if name[0] == '$' {
		return resolver(strings.TrimPrefix(name, "$"))
	}
	return name, nil
}

func (r *Runner) resolvePath(urlPath string) (string, error) {
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

func (r *Runner) Login() error {
	loginUrl, err := r.resolvePath(r.Definition.Login.Path)
	if err != nil {
		return err
	}

	r.logger.WithField("url", loginUrl).Info("Attempting to login")

	err = r.Browser.Open(loginUrl)
	if err != nil {
		return err
	}

	r.logger.
		WithFields(log.Fields{"code": r.Browser.StatusCode(), "page": r.Browser.Url()}).
		Debugf("Finished request")

	fm, err := r.Browser.Form(r.Definition.Login.FormSelector)
	if err != nil {
		return err
	}

	for name, val := range r.Definition.Login.Inputs {
		r.logger.
			WithFields(log.Fields{"key": name, "form": r.Definition.Login.FormSelector, "val": val}).
			Debugf("Filling input of form")

		resolved, err := r.resolveVariable(val, func(name string) (string, error) {
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

	r.logger.Debug("Submitting login form")

	if err = fm.Submit(); err != nil {
		r.logger.WithError(err).Error("Login failed")
		return err
	}

	r.logger.
		WithFields(log.Fields{"code": r.Browser.StatusCode(), "page": r.Browser.Url()}).
		Debugf("Finished request")

	if err = r.Definition.Login.hasError(r.Browser); err != nil {
		r.logger.WithError(err).Info("Failed to login")
		return err
	}

	r.logger.Info("Successfully logged in")
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
		r.logger.Info("Testing search mode %s", mode.Key)
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
	searchUrl, err := r.resolvePath(r.Definition.Search.Path)
	if err != nil {
		return nil, err
	}

	r.logger.
		WithFields(log.Fields{"page": searchUrl, "query": query}).
		Infof("Opening search page")

	err = r.Browser.Open(searchUrl)
	if err != nil {
		return nil, err
	}

	r.logger.
		WithFields(log.Fields{"code": r.Browser.StatusCode(), "page": r.Browser.Url()}).
		Debugf("Finished request")

	vals := url.Values{}

	for name, val := range r.Definition.Search.Inputs {
		resolved, err := r.resolveVariable(val, func(name string) (string, error) {
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

	r.logger.
		WithFields(log.Fields{"params": vals, "page": searchUrl}).
		Debugf("Submitting page with form params")

	err = r.Browser.OpenForm(searchUrl, vals)
	if err != nil {
		return nil, err
	}

	r.logger.
		WithFields(log.Fields{"code": r.Browser.StatusCode(), "page": r.Browser.Url()}).
		Debugf("Finished request")

	items := []torznab.ResultItem{}
	timer := time.Now()
	rows := r.Browser.Find(r.Definition.Search.Rows.Selector)

	r.logger.
		WithFields(log.Fields{"rows": rows.Length(), "selector": r.Definition.Search.Rows.Selector}).
		Debugf("Found %d rows", rows.Length())

	for i := 0; i < rows.Length(); i++ {
		row := map[string]string{}

		for field, block := range r.Definition.Search.Fields {
			r.logger.
				WithFields(log.Fields{"row": i + 1, "block": block}).
				Debugf("Processing field %q", field)

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

		r.logger.
			WithFields(log.Fields{"row": i + 1, "data": row}).
			Debugf("Finished row %d", i+1)

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

	r.logger.WithFields(log.Fields{"time": time.Now().Sub(timer)}).Infof("Query returned %d results", len(items))
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
