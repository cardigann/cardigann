package indexer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
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
	bow.SetAttribute(browser.MetaRefreshHandling, true)

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

// this should eventually upstream into surf browser
func (r *Runner) handleMetaRefreshHeader() error {
	h := r.Browser.ResponseHeaders()

	if refresh := h.Get("Refresh"); refresh != "" {
		if s := regexp.MustCompile(`\s*;\s*`).Split(refresh, 2); len(s) == 2 {
			r.Logger.
				WithField("fields", s).
				Debug("Found refresh header")

			u, err := r.resolvePath(strings.TrimPrefix(s[1], "url="))
			if err != nil {
				return err
			}

			return r.openPage(u)
		}
	}
	return nil
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

	if err = r.handleMetaRefreshHeader(); err != nil {
		return err
	}

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

func (r *Runner) isLoginRequired() bool {
	if r.Definition.Login.Test.Path == "" {
		return true
	}

	r.Logger.
		WithField("path", r.Definition.Login.Test).
		Debug("Testing if login is needed")

	testUrl, err := r.resolvePath(r.Definition.Login.Test.Path)
	if err != nil {
		r.Logger.WithError(err).Warn("Failed to resolve path")
		return true
	}

	err = r.openPage(testUrl)
	if err != nil {
		r.Logger.WithError(err).Warn("Failed to open page")
		return true
	}

	if testUrl == r.Browser.Url().String() {
		r.Logger.Debug("No login needed, already logged in")
		return false
	}

	r.Logger.Debug("Login is required")
	return true
}

func (r *Runner) login() error {
	filterLogger = r.Logger

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

	if r.isLoginRequired() {
		if err := r.login(); err != nil {
			r.Logger.WithError(err).Error("Login failed")
			return nil, err
		}
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
	if queryCatIDs, ok := query["cat"].([]int); ok {
		queryCats := torznab.AllCategories.Subset(queryCatIDs...)

		// resolve query categories to the exact local, or the local based on parent cat
		for _, id := range r.Capabilities().Categories.ResolveAll(queryCats...) {
			localCats = append(localCats, id)
		}

		r.Logger.
			WithFields(logrus.Fields{"querycats": queryCatIDs, "local": localCats}).
			Debugf("Resolved torznab cats to local")
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
				Debugf("Processed $raw input")

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
	limit, hasLimit := query.Limit()

	r.Logger.
		WithFields(logrus.Fields{
		"rows":     rows.Length(),
		"selector": r.Definition.Search.Rows.Selector,
		"limit":    limit,
	}).Debugf("Found %d rows", rows.Length())

	items := []torznab.ResultItem{}

	for i := 0; i < rows.Length(); i++ {
		if hasLimit && len(items) >= limit {
			break
		}

		item, err := r.extractItem(i+1, rows.Eq(i))
		if err != nil {
			return nil, err
		}

		var matchCat bool
		if len(localCats) > 0 {
			for _, catId := range localCats {
				if catId == item.Category {
					matchCat = true
				}
			}

			if !matchCat {
				r.Logger.
					WithFields(logrus.Fields{"id": item.Category, "localCats": localCats}).
					Debug("Skipping non-matching category")
				continue
			}
		}

		if mappedCat, ok := r.Definition.Capabilities.Categories[item.Category]; ok {
			item.Category = mappedCat.ID
		} else {
			item.Category = item.Category + torznab.CustomCategoryOffset
		}

		items = append(items, item)
	}

	r.Logger.
		WithFields(logrus.Fields{"time": time.Now().Sub(timer)}).
		Infof("Query returned %d results", len(items))

	return items, nil
}

func (r *Runner) extractItem(rowIdx int, selection *goquery.Selection) (torznab.ResultItem, error) {
	row := map[string]string{}

	html, _ := goquery.OuterHtml(selection)
	r.Logger.WithFields(logrus.Fields{"html": gohtml.Format(html)}).Debug("Processing row")

	for _, item := range r.Definition.Search.Fields {
		r.Logger.
			WithFields(logrus.Fields{"row": rowIdx, "block": item.Block.String()}).
			Debugf("Processing field %q", item.Field)

		val, err := item.Block.MatchText(selection)
		if err != nil {
			return torznab.ResultItem{}, err
		}

		r.Logger.
			WithFields(logrus.Fields{"row": rowIdx, "output": val}).
			Debugf("Finished processing field %q", item.Field)

		row[item.Field] = val
	}

	item := torznab.ResultItem{
		Site:            r.Definition.Site,
		MinimumRatio:    1,
		MinimumSeedTime: time.Hour * 48,
	}

	r.Logger.
		WithFields(logrus.Fields{"row": rowIdx, "data": row}).
		Debugf("Finished row %d", rowIdx)

	for key, val := range row {
		switch key {
		case "download":
			u, err := r.resolvePath(val)
			if err != nil {
				r.Logger.Warnf("Search result row #%d has malformed url in %s", rowIdx, key)
				continue
			}
			item.Link = u
		case "details":
			u, err := r.resolvePath(val)
			if err != nil {
				r.Logger.Warnf("Search result row #%d has malformed url in %s", rowIdx, key)
				continue
			}
			item.GUID = u
		case "comments":
			u, err := r.resolvePath(val)
			if err != nil {
				r.Logger.Warnf("Search result row #%d has malformed url in %s", rowIdx, key)
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
				r.Logger.Warnf("Search result row #%d has malformed categoryid: %s", rowIdx, err.Error())
				continue
			}
			item.Category = catID
		case "size":
			bytes, err := humanize.ParseBytes(val)
			if err != nil {
				r.Logger.Warnf("Search result row #%d has malformed size: %s", rowIdx, err.Error())
				continue
			}
			item.Size = bytes
		case "leechers":
			leechers, err := strconv.Atoi(val)
			if err != nil {
				r.Logger.Warnf("Search result row #%d has malformed leechers value in %s", rowIdx, key)
				continue
			}
			item.Peers += leechers
		case "seeders":
			seeders, err := strconv.Atoi(val)
			if err != nil {
				r.Logger.Warnf("Search result row #%d has malformed seeders value in %s", rowIdx, key)
				continue
			}
			item.Seeders = seeders
			item.Peers += seeders
		case "date":
			t, err := time.Parse(filterTimeFormat, val)
			if err != nil {
				r.Logger.Warnf("Search result row #%d has malformed time value in %s", rowIdx, key)
				continue
			}
			item.PublishDate = t
		default:
			r.Logger.Warnf("Search result row #%d has unknown field %s", rowIdx, key)
			continue
		}
	}

	if item.GUID == "" && item.Link != "" {
		item.GUID = item.Link
	}

	if r.hasDateHeader() {
		date, err := r.extractDateHeader(selection)
		if err != nil {
			return torznab.ResultItem{}, err
		}

		item.PublishDate = date
	}

	return item, nil
}

func (r *Runner) hasDateHeader() bool {
	return !r.Definition.Search.Rows.DateHeaders.IsEmpty()
}

func (r *Runner) extractDateHeader(selection *goquery.Selection) (time.Time, error) {
	dateHeaders := r.Definition.Search.Rows.DateHeaders

	r.Logger.
		WithFields(logrus.Fields{"selector": dateHeaders.String()}).
		Debugf("Searching for date header")

	prev := selection.PrevAllFiltered(dateHeaders.Selector).First()
	if prev.Length() == 0 {
		return time.Time{}, fmt.Errorf("No date header row found")
	}

	dv, _ := dateHeaders.Text(prev.First())
	return time.Parse(filterTimeFormat, dv)
}

func (r *Runner) Download(u string) (io.ReadCloser, http.Header, error) {
	if r.isLoginRequired() {
		if err := r.login(); err != nil {
			r.Logger.WithError(err).Error("Login failed")
			return nil, http.Header{}, err
		}
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
