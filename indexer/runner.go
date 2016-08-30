package indexer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
	"github.com/cardigann/cardigann/config"
	"github.com/cardigann/cardigann/torznab"
	"github.com/dustin/go-humanize"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"github.com/yosssi/gohtml"
)

var (
	_ torznab.Indexer = &Runner{}
)

type Runner struct {
	Definition *IndexerDefinition
	Browser    browser.Browsable
	Config     config.Config
	Logger     logrus.FieldLogger
	caps       torznab.Capabilities
}

func NewRunner(def *IndexerDefinition, conf config.Config) *Runner {
	bow := surf.NewBrowser()
	bow.SetUserAgent(agent.Chrome())
	bow.SetAttribute(browser.SendReferer, false)
	bow.SetAttribute(browser.MetaRefreshHandling, false)

	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	return &Runner{
		Definition: def,
		Browser:    bow,
		Config:     conf,
		Logger:     logger.WithFields(logrus.Fields{"site": def.Site}),
	}
}

func (r *Runner) applyTemplate(name, tpl string, ctx interface{}) (string, error) {
	tmpl, err := template.New(name).Parse(tpl)
	if err != nil {
		return "", err
	}
	b := &bytes.Buffer{}
	err = tmpl.Execute(b, ctx)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (r *Runner) currentURL() (*url.URL, error) {
	if u := r.Browser.Url(); u != nil {
		return u, nil
	}

	if configURL, ok, _ := r.Config.Get(r.Definition.Site, "url"); ok {
		return url.Parse(configURL)
	}

	return url.Parse(r.Definition.Links[0])
}

func (r *Runner) resolvePath(urlPath string) (string, error) {
	base, err := r.currentURL()
	if err != nil {
		return "", err
	}

	u, err := url.Parse(urlPath)
	if err != nil {
		log.Fatal(err)
	}

	resolved := base.ResolveReference(u)
	r.Logger.
		WithFields(logrus.Fields{"base": base.String(), "u": resolved.String()}).
		Debugf("Resolving url")

	return resolved.String(), nil
}

func (r *Runner) openPage(u string) error {
	r.Logger.WithField("url", u).Debug("Attempting to open page")

	err := r.Browser.Open(u)
	if err != nil {
		return err
	}

	r.Logger.
		WithFields(logrus.Fields{"code": r.Browser.StatusCode(), "page": r.Browser.Url()}).
		Debugf("Finished request")

	return r.cachePage()
}

func (r *Runner) cachePage() error {
	tmpfile, err := ioutil.TempFile("", r.Definition.Site)
	if err != nil {
		return err
	}

	body := strings.NewReader(r.Browser.Body())
	io.Copy(tmpfile, body)
	defer tmpfile.Close()

	r.Logger.
		WithFields(logrus.Fields{"file": "file://" + tmpfile.Name()}).
		Debugf("Wrote page output to cache")

	return nil
}

func (r *Runner) fillAndSubmitForm(loginURL, formSelector string, vals map[string]string) error {
	r.Logger.
		WithFields(logrus.Fields{"url": loginURL, "form": formSelector, "vals": vals}).
		Debugf("Filling and submitting login form")

	if err := r.openPage(loginURL); err != nil {
		return err
	}

	fm, err := r.Browser.Form(formSelector)
	if err != nil {
		return err
	}

	for name, value := range vals {
		r.Logger.
			WithFields(logrus.Fields{"key": name, "form": formSelector, "val": value}).
			Debugf("Filling input of form")

		if err = fm.Input(name, value); err != nil {
			r.Logger.WithError(err).Error("Filling input failed")
			return err
		}
	}

	r.Logger.Debug("Submitting login form")
	defer r.cachePage()

	if err = fm.Submit(); err != nil {
		r.Logger.WithError(err).Error("Login failed")
		return err
	}

	r.Logger.
		WithFields(logrus.Fields{"code": r.Browser.StatusCode(), "page": r.Browser.Url()}).
		Debugf("Submitted login form")

	return nil
}

func (r *Runner) postForm(loginURL string, vals map[string]string) error {
	r.Logger.
		WithFields(logrus.Fields{"url": loginURL, "vals": vals}).
		Debugf("Posting login form")

	data := url.Values{}
	for key, value := range vals {
		data.Add(key, value)
	}

	defer r.cachePage()

	if err := r.Browser.PostForm(loginURL, data); err != nil {
		r.Logger.WithError(err).Error("Login failed")
		return err
	}

	r.Logger.
		WithFields(logrus.Fields{"code": r.Browser.StatusCode(), "page": r.Browser.Url()}).
		Debugf("Posted login page")

	return nil
}

func (r *Runner) extractInputLogins() (map[string]string, error) {
	result := map[string]string{}

	cfg, err := r.Config.Section(r.Definition.Site)
	if err != nil {
		return nil, err
	}

	ctx := struct {
		Config map[string]string
	}{
		cfg,
	}

	for name, val := range r.Definition.Login.Inputs {
		resolved, err := r.applyTemplate("login_inputs", val, ctx)
		if err != nil {
			return nil, err
		}

		r.Logger.
			WithFields(logrus.Fields{"key": name, "val": resolved}).
			Debugf("Resolved login input template")

		result[name] = resolved
	}

	return result, nil
}

func (r *Runner) login() error {
	filterLogger = r.Logger
	filterCategoryMapping = r.Capabilities().Categories

	loginUrl, err := r.resolvePath(r.Definition.Login.Path)
	if err != nil {
		return err
	}

	vals, err := r.extractInputLogins()
	if err != nil {
		return err
	}

	switch r.Definition.Login.Method {
	case "", loginMethodForm:
		if err = r.fillAndSubmitForm(loginUrl, r.Definition.Login.FormSelector, vals); err != nil {
			return err
		}
	case loginMethodPost:
		if err = r.postForm(loginUrl, vals); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown login method %q", r.Definition.Login.Method)
	}

	if err = r.Definition.Login.hasError(r.Browser); err != nil {
		r.Logger.WithError(err).Error("Failed to login")
		return err
	}

	r.Logger.Info("Successfully logged in")
	return nil
}

func (r *Runner) Info() torznab.Info {
	return torznab.Info{
		ID:       r.Definition.Site,
		Title:    r.Definition.Name,
		Language: r.Definition.Language,
		Link:     r.Definition.Links[0],
	}
}

func (r *Runner) Test() error {
	filterLogger = r.Logger

	for _, mode := range r.Capabilities().SearchModes {
		query := torznab.Query{
			"t":     mode.Key,
			"limit": 5,
		}

		switch mode.Key {
		case "tv-search":
			query["cat"] = []int{
				torznab.CategoryTV.ID,
				torznab.CategoryTV_HD.ID,
				torznab.CategoryTV_SD.ID,
			}
		}

		r.Logger.Infof("Testing search mode %q", mode.Key)
		results, err := r.Search(query)
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
		}
	}

	return nil
}

func (r *Runner) Capabilities() torznab.Capabilities {
	return torznab.Capabilities(r.Definition.Capabilities)
}

func (r *Runner) Search(query torznab.Query) ([]torznab.ResultItem, error) {
	filterLogger = r.Logger
	filterCategoryMapping = r.Capabilities().Categories

	if err := r.login(); err != nil {
		r.Logger.WithError(err).Error("Login failed")
		return nil, err
	}

	searchUrl, err := r.resolvePath(r.Definition.Search.Path)
	if err != nil {
		return nil, err
	}

	r.Logger.
		WithFields(logrus.Fields{"query": query}).
		Infof("Searching indexer")

	if err := r.openPage(searchUrl); err != nil {
		return nil, err
	}

	localCats := []int{}

	if unmappedCats, ok := query["cat"].([]int); ok {
		localCats = r.Capabilities().Categories.ReverseMap(unmappedCats)
	}

	inputCtx := struct {
		Query      torznab.Query
		Keywords   string
		Categories []int
	}{
		query,
		query.Keywords(),
		localCats,
	}

	vals := url.Values{}

	for name, val := range r.Definition.Search.Inputs {
		resolved, err := r.applyTemplate("search_inputs", val, inputCtx)
		if err != nil {
			return nil, err
		}
		switch name {
		case "$raw":
			parsedVals, err := url.ParseQuery(resolved)
			if err != nil {
				return nil, fmt.Errorf("Error parsing $raw input: %s", err.Error())
			}

			r.Logger.
				WithFields(logrus.Fields{"source": val, "parsed": parsedVals}).
				Infof("Processed $raw input")

			for k, values := range parsedVals {
				for _, val := range values {
					vals.Add(k, val)
				}
			}
		default:
			vals.Add(name, resolved)
		}
	}

	r.Logger.
		WithFields(logrus.Fields{"params": vals, "page": searchUrl}).
		Debugf("Submitting page with form params")

	err = r.Browser.OpenForm(searchUrl, vals)
	if err != nil {
		return nil, err
	}

	r.Logger.
		WithFields(logrus.Fields{"code": r.Browser.StatusCode(), "page": r.Browser.Url()}).
		Debugf("Finished opening form")

	items := []torznab.ResultItem{}
	timer := time.Now()
	dom := r.Browser.Dom()

	// merge following rows for After selector
	if after := r.Definition.Search.Rows.After; after > 0 {
		rows := dom.Find(r.Definition.Search.Rows.Selector)
		for i := 0; i < rows.Length(); i += 1 + after {
			rows.Eq(i).AppendSelection(rows.Slice(i+1, i+1+after).Find("td"))
			rows.Slice(i+1, i+1+after).Remove()
		}
	}

	rows := dom.Find(r.Definition.Search.Rows.Selector)
	limit, hasLimit := query["limit"].(int)

	r.Logger.
		WithFields(logrus.Fields{
		"rows":     rows.Length(),
		"selector": r.Definition.Search.Rows.Selector,
		"limit":    limit,
	}).
		Debugf("Found %d rows", rows.Length())

	for i := 0; i < rows.Length() && (!hasLimit || len(items) < limit); i++ {
		row := map[string]string{}

		html, _ := goquery.OuterHtml(rows.Eq(i))
		r.Logger.WithFields(logrus.Fields{"html": gohtml.Format(html)}).Debug("Processing row")

		for _, item := range r.Definition.Search.Fields {
			r.Logger.
				WithFields(logrus.Fields{"row": i + 1, "block": item.Block.String()}).
				Debugf("Processing field %q", item.Field)

			val, err := item.Block.MatchText(rows.Eq(i))
			if err != nil {
				return nil, err
			}

			r.Logger.
				WithFields(logrus.Fields{"row": i + 1, "output": val}).
				Debugf("Finished processing field %q", item.Field)

			row[item.Field] = val
		}

		item := torznab.ResultItem{
			Site:            r.Definition.Site,
			MinimumRatio:    1,
			MinimumSeedTime: time.Hour * 48,
		}

		r.Logger.
			WithFields(logrus.Fields{"row": i + 1, "data": row}).
			Debugf("Finished row %d", i+1)

		for key, val := range row {
			switch key {
			case "download":
				u, err := r.resolvePath(val)
				if err != nil {
					r.Logger.Warnf("Search result row #%d has malformed url in %s", i+1, key)
					continue
				}
				item.Link = u
			case "details":
				u, err := r.resolvePath(val)
				if err != nil {
					r.Logger.Warnf("Search result row #%d has malformed url in %s", i+1, key)
					continue
				}
				item.GUID = u
			case "comments":
				u, err := r.resolvePath(val)
				if err != nil {
					r.Logger.Warnf("Search result row #%d has malformed url in %s", i+1, key)
					continue
				}
				item.Comments = u
			case "title":
				item.Title = val
			case "description":
				item.Description = val
			case "category":
				catID, err := strconv.Atoi(val)
				if err != nil {
					r.Logger.Warnf("Search result row #%d has malformed categoryid: %s", i+1, err.Error())
					continue
				}
				item.Category = catID
			case "size":
				bytes, err := humanize.ParseBytes(val)
				if err != nil {
					r.Logger.Warnf("Search result row #%d has malformed size: %s", i+1, err.Error())
					continue
				}
				item.Size = bytes
			case "leechers":
				leechers, err := strconv.Atoi(val)
				if err != nil {
					r.Logger.Warnf("Search result row #%d has malformed leechers value in %s", i+1, key)
					continue
				}
				item.Peers += leechers
			case "seeders":
				seeders, err := strconv.Atoi(val)
				if err != nil {
					r.Logger.Warnf("Search result row #%d has malformed seeders value in %s", i+1, key)
					continue
				}
				item.Seeders = seeders
				item.Peers += seeders
			case "date":
				t, err := time.Parse(filterTimeFormat, val)
				if err != nil {
					r.Logger.Warnf("Search result row #%d has malformed time value in %s", i+1, key)
					continue
				}
				item.PublishDate = t
			default:
				return nil, fmt.Errorf("Unknown field %q", key)
			}
		}

		skipItem := false

		// some trackers have empty rows when there are no results
		if item.Title == "" {
			return nil, nil
		}

		if item.GUID == "" && item.Link != "" {
			item.GUID = item.Link
		}

		// if there is a dateheaders field, we need to look for a preceeding date header
		if dateHeaders := r.Definition.Search.Rows.DateHeaders; !dateHeaders.IsEmpty() {
			r.Logger.
				WithFields(logrus.Fields{"selector": dateHeaders.String()}).
				Debugf("Searching for date header")

			prev := rows.Eq(i).PrevAllFiltered(dateHeaders.Selector).First()
			if prev.Length() == 0 {
				r.Logger.
					WithFields(logrus.Fields{"row": i + 1}).
					Warnf("No preceding date header found for row %#d", i+1)
			}

			dv, _ := dateHeaders.Text(prev.First())
			date, err := time.Parse(filterTimeFormat, dv)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse date header: %s", err.Error())
			}

			item.PublishDate = date
		}

		// some trackers don't support filtering by categories, so do it for them
		if catFilters, hasCats := query["cat"].([]int); hasCats {
			var catMatch bool
			for _, catId := range catFilters {
				r.Logger.Debugf("Checking item cat %d against query cat %d", item.Category, catId)
				if catId == item.Category {
					catMatch = true
				}
			}
			if !catMatch {
				r.Logger.Debugf("Skipping row due to non-matching category")
				skipItem = skipItem || !catMatch
			}
		}

		if !skipItem {
			items = append(items, item)
		}
	}

	r.Logger.
		WithFields(logrus.Fields{"time": time.Now().Sub(timer)}).
		Infof("Query returned %d results", len(items))

	return items, nil
}

func (r *Runner) Download(u string) (io.ReadCloser, http.Header, error) {
	if err := r.login(); err != nil {
		r.Logger.WithError(err).Error("Login failed")
		return nil, http.Header{}, err
	}

	fullUrl, err := r.resolvePath(u)
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
