package indexer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
	"github.com/cardigann/cardigann/config"
	"github.com/cardigann/cardigann/logger"
	"github.com/cardigann/cardigann/torznab"
	"github.com/dustin/go-humanize"
	"github.com/f2prateek/train"
	trainlog "github.com/f2prateek/train/log"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"github.com/headzoo/surf/jar"
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

	if os.Getenv("DEBUG_HTTP") != "" {
		bow.SetTransport(train.Transport(trainlog.New(os.Stderr, trainlog.Basic)))
	}

	return &Runner{
		Definition: def,
		Browser:    bow,
		Config:     conf,
		Logger:     logger.Logger.WithFields(logrus.Fields{"site": def.Site}),
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
	if strings.Contains(tpl, "{{") {
		r.Logger.
			WithFields(logrus.Fields{"src": tpl, "result": b.String()}).
			Debugf("Processed template")
	}
	return b.String(), nil
}

func (r *Runner) currentURL() (*url.URL, error) {
	if u := r.Browser.Url(); u != nil {
		return u, nil
	}

	configURL, ok, _ := r.Config.Get(r.Definition.Site, "url")
	if ok && r.testURLWorks(configURL) {
		return url.Parse(configURL)
	}

	for _, u := range r.Definition.Links {
		if u != configURL && r.testURLWorks(u) {
			return url.Parse(u)
		}
	}

	return nil, errors.New("No working urls found")
}

func (r *Runner) testURLWorks(u string) bool {
	r.Logger.WithField("url", u).Debugf("Checking connectivity to url")

	err := r.Browser.Open(u)
	if err != nil {
		r.Logger.WithError(err).Warn("URL check failed")
		return false
	} else if r.Browser.StatusCode() != http.StatusOK {
		r.Logger.Warn("URL returned non-ok status")
		return false
	}

	return true
}

func (r *Runner) resolvePath(urlPath string) (string, error) {
	base, err := r.currentURL()
	if err != nil {
		return "", err
	}

	u, err := url.Parse(urlPath)
	if err != nil {
		return "", err
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

func parseCookieString(cookie string) []*http.Cookie {
	h := http.Header{"Cookie": []string{cookie}}
	r := http.Request{Header: h}
	return r.Cookies()
}

func (r *Runner) loginWithCookie(loginURL string, cookie string) error {
	u, err := url.Parse(loginURL)
	if err != nil {
		return err
	}

	cookies := parseCookieString(cookie)

	r.Logger.
		WithFields(logrus.Fields{"url": loginURL, "cookies": cookies}).
		Debugf("Setting cookies for login")

	cj := jar.NewMemoryCookies()
	cj.SetCookies(u, cookies)

	r.Browser.SetCookieJar(cj)
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

func (r *Runner) isLoginRequired() (bool, error) {
	if r.Definition.Login.Test.Path == "" {
		return true, nil
	}

	r.Logger.
		WithField("path", r.Definition.Login.Test).
		Debug("Testing if login is needed")

	testUrl, err := r.resolvePath(r.Definition.Login.Test.Path)
	if err != nil {
		return true, err
	}

	err = r.openPage(testUrl)
	if err != nil {
		r.Logger.WithError(err).Warn("Failed to open page")
		return true, nil
	}

	if testUrl == r.Browser.Url().String() {
		r.Logger.Debug("No login needed, already logged in")
		return false, nil
	}

	r.Logger.Debug("Login is required")
	return true, nil
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
	case loginMethodCookie:
		if err = r.loginWithCookie(loginUrl, vals["cookie"]); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown login method %q", r.Definition.Login.Method)
	}

	if len(r.Definition.Login.Error) > 0 {
		r.Logger.
			WithField("block", r.Definition.Login.Error).
			Debug("Testing if login succeeded")

		if err = r.Definition.Login.hasError(r.Browser); err != nil {
			r.Logger.WithError(err).Error("Failed to login")
			return err
		}
	}

	if r.Definition.Login.Test.Path != "" {
		if required, err := r.isLoginRequired(); err != nil {
			return err
		} else if required {
			return errors.New("Login check after login failed")
		}
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

func (r *Runner) Capabilities() torznab.Capabilities {
	return r.Definition.Capabilities.ToTorznab()
}

type extractedItem struct {
	torznab.ResultItem
	LocalCategoryID string
}

func (r *Runner) Search(query torznab.Query) ([]torznab.ResultItem, error) {
	filterLogger = r.Logger

	if required, err := r.isLoginRequired(); err != nil {
		return nil, err
	} else if required {
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

	localCats := []string{}
	if queryCatIDs, ok := query["cat"].([]int); ok {
		queryCats := torznab.AllCategories.Subset(queryCatIDs...)

		// resolve query categories to the exact local, or the local based on parent cat
		for _, id := range r.Definition.Capabilities.CategoryMap.ResolveAll(queryCats...) {
			localCats = append(localCats, id)
		}

		r.Logger.
			WithFields(logrus.Fields{"querycats": queryCatIDs, "local": localCats}).
			Debugf("Resolved torznab cats to local")
	}

	inputCtx := struct {
		Query      torznab.Query
		Keywords   string
		Categories []string
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
				r.Logger.WithError(err).Warn(err)
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

	// apply Remove if it exists
	if remove := r.Definition.Search.Rows.Remove; remove != "" {
		matching := dom.Find(r.Definition.Search.Rows.Selector).Filter(remove)
		r.Logger.
			WithFields(logrus.Fields{"selector": remove}).
			Debugf("Applying remove to %d rows", matching.Length())
		matching.Remove()
	}

	rows := dom.Find(r.Definition.Search.Rows.Selector)
	limit, hasLimit := query.Limit()

	r.Logger.
		WithFields(logrus.Fields{
		"rows":     rows.Length(),
		"selector": r.Definition.Search.Rows.Selector,
		"limit":    limit,
	}).Debugf("Found %d rows", rows.Length())

	extracted := []extractedItem{}

	for i := 0; i < rows.Length(); i++ {
		if hasLimit && len(extracted) >= limit {
			break
		}

		item, err := r.extractItem(i+1, rows.Eq(i))
		if err != nil {
			return nil, err
		}

		var matchCat bool
		if len(localCats) > 0 {
			for _, catId := range localCats {
				if catId == item.LocalCategoryID {
					matchCat = true
				}
			}

			if !matchCat {
				r.Logger.
					WithFields(logrus.Fields{"id": item.LocalCategoryID, "localCats": localCats}).
					Debug("Skipping non-matching category")
				continue
			}
		}

		if mappedCat, ok := r.Definition.Capabilities.CategoryMap[item.LocalCategoryID]; ok {
			item.Category = mappedCat.ID
		} else {
			r.Logger.
				WithFields(logrus.Fields{"localId": item.LocalCategoryID}).
				Warn("Unknown local category")

			if intCatId, err := strconv.Atoi(item.LocalCategoryID); err == nil {
				item.Category = intCatId + torznab.CustomCategoryOffset
			}
		}

		extracted = append(extracted, item)
	}

	r.Logger.
		WithFields(logrus.Fields{"time": time.Now().Sub(timer)}).
		Infof("Query returned %d results", len(extracted))

	items := []torznab.ResultItem{}
	for _, item := range extracted {
		items = append(items, item.ResultItem)
	}

	return items, nil
}

func (r *Runner) extractItem(rowIdx int, selection *goquery.Selection) (extractedItem, error) {
	row := map[string]string{}

	html, _ := goquery.OuterHtml(selection)
	r.Logger.WithFields(logrus.Fields{"html": gohtml.Format(html)}).Debug("Processing row")

	for _, item := range r.Definition.Search.Fields {
		r.Logger.
			WithFields(logrus.Fields{"row": rowIdx, "block": item.Block.String()}).
			Debugf("Processing field %q", item.Field)

		val, err := item.Block.MatchText(selection)
		if err != nil {
			return extractedItem{}, err
		}

		r.Logger.
			WithFields(logrus.Fields{"row": rowIdx, "output": val}).
			Debugf("Finished processing field %q", item.Field)

		row[item.Field] = val
	}

	item := extractedItem{
		ResultItem: torznab.ResultItem{
			Site:            r.Definition.Site,
			MinimumRatio:    1,
			MinimumSeedTime: time.Hour * 48,
		},
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
			item.LocalCategoryID = val
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
			return extractedItem{}, err
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
	if required, err := r.isLoginRequired(); required {
		if err := r.login(); err != nil {
			r.Logger.WithError(err).Error("Login failed")
			return nil, http.Header{}, err
		}
	} else if err != nil {
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
