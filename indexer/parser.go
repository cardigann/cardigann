package indexer

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cardigann/cardigann/torznab"
	"github.com/headzoo/surf/browser"

	"gopkg.in/yaml.v2"
)

type IndexerDefinition struct {
	Site         string                 `yaml:"site"`
	Settings     []settingsField        `yaml:"settings"`
	Name         string                 `yaml:"name"`
	Description  string                 `yaml:"description"`
	Language     string                 `yaml:"language"`
	Links        stringorslice          `yaml:"links"`
	Capabilities capabilitiesBlock      `yaml:"caps"`
	Login        loginBlock             `yaml:"login"`
	Ratio        ratioBlock             `yaml:"ratio"`
	Search       searchBlock            `yaml:"search"`
	stats        IndexerDefinitionStats `yaml:"-"`
}

type IndexerDefinitionStats struct {
	Size    int64
	ModTime time.Time
	Hash    string
	Source  string
}

func (id *IndexerDefinition) Stats() IndexerDefinitionStats {
	return id.stats
}

type settingsField struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Label string `yaml:"label"`
}

func ParseDefinitionFile(f *os.File) (*IndexerDefinition, error) {
	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, err
	}

	def, err := ParseDefinition(b)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	def.stats.ModTime = fi.ModTime()
	return def, err
}

func ParseDefinition(src []byte) (*IndexerDefinition, error) {
	def := IndexerDefinition{
		Language:     "en-us",
		Capabilities: capabilitiesBlock{},
		Login: loginBlock{
			FormSelector: "form",
			Inputs:       inputsBlock{},
		},
		Search: searchBlock{},
	}

	if err := yaml.Unmarshal(src, &def); err != nil {
		return nil, err
	}

	if len(def.Settings) == 0 {
		def.Settings = defaultSettingsFields()
	}

	def.stats = IndexerDefinitionStats{
		Size:    int64(len(src)),
		ModTime: time.Now(),
		Hash:    fmt.Sprintf("%x", sha1.Sum(src)),
	}

	return &def, nil
}

func defaultSettingsFields() []settingsField {
	return []settingsField{
		{Name: "username", Label: "Username", Type: "text"},
		{Name: "password", Label: "Password", Type: "password"},
	}
}

type inputsBlock map[string]string

type errorBlockOrSlice []errorBlock

// UnmarshalYAML implements the Unmarshaller interface.
func (e *errorBlockOrSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var blockType errorBlock
	if err := unmarshal(&blockType); err == nil {
		*e = errorBlockOrSlice{blockType}
		return nil
	}

	var sliceType []errorBlock
	if err := unmarshal(&sliceType); err == nil {
		*e = errorBlockOrSlice(sliceType)
		return nil
	}

	return errors.New("Failed to unmarshal errorBlockOrSlice")
}

type errorBlock struct {
	Path     string        `yaml:"path"`
	Selector string        `yaml:"selector"`
	Message  selectorBlock `yaml:"message"`
}

func (e *errorBlock) matchPage(browser browser.Browsable) bool {
	if e.Path != "" {
		return e.Path == browser.Url().Path
	} else if e.Selector != "" {
		return browser.Find(e.Selector).Length() > 0
	}
	return false
}

func (e *errorBlock) errorText(from *goquery.Selection) (string, error) {
	if !e.Message.IsEmpty() {
		return e.Message.MatchText(from)
	} else if e.Selector != "" {
		return from.Find(e.Selector).Text(), nil
	}
	return "", errors.New("Error declaration must have either Message block or Selection")
}

type pageTestBlock struct {
	Path     string `yaml:"path"`
	Selector string `yaml:"selector"`
}

func (t *pageTestBlock) IsEmpty() bool {
	return t.Path == "" && t.Selector == ""
}

const (
	loginMethodPost   = "post"
	loginMethodGet    = "get"
	loginMethodForm   = "form"
	loginMethodCookie = "cookie"
)

type loginBlock struct {
	Path         string            `yaml:"path"`
	FormSelector string            `yaml:"form"`
	Method       string            `yaml:"method"`
	Inputs       inputsBlock       `yaml:"inputs,omitempty"`
	Error        errorBlockOrSlice `yaml:"error,omitempty"`
	Test         pageTestBlock     `yaml:"test,omitempty"`
}

func (l *loginBlock) IsEmpty() bool {
	return l.Path == "" && l.Method == ""
}

func (l *loginBlock) hasError(browser browser.Browsable) error {
	for _, e := range l.Error {
		if e.matchPage(browser) {
			msg, err := e.errorText(browser.Dom())
			if err != nil {
				return err
			}
			return errors.New(strings.TrimSpace(msg))
		}
	}

	return nil
}

type fieldBlock struct {
	Field string
	Block selectorBlock
}

type fieldsListBlock []fieldBlock

func (f *fieldsListBlock) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Unmarshal as a MapSlice to preserve order of fields
	var fields yaml.MapSlice
	if err := unmarshal(&fields); err != nil {
		return errors.New("Failed to unmarshal fieldsListBlock")
	}

	// FIXME: there has got to be a better way to do this
	for _, item := range fields {
		b, err := yaml.Marshal(item.Value)
		if err != nil {
			return err
		}
		var sb selectorBlock
		if err = yaml.Unmarshal(b, &sb); err != nil {
			return err
		}
		*f = append(*f, fieldBlock{
			Field: item.Key.(string),
			Block: sb,
		})
	}

	return nil
}

type rowsBlock struct {
	selectorBlock
	After       int           `yaml:"after"`
	Remove      string        `yaml:"remove"`
	DateHeaders selectorBlock `yaml:"dateheaders"`
}

func (r *rowsBlock) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var sb selectorBlock
	if err := unmarshal(&sb); err != nil {
		return errors.New("Failed to unmarshal rowsBlock")
	}

	var rb struct {
		After       int           `yaml:"after"`
		Remove      string        `yaml:"remove"`
		DateHeaders selectorBlock `yaml:"dateheaders"`
	}
	if err := unmarshal(&rb); err != nil {
		return errors.New("Failed to unmarshal rowsBlock")
	}

	r.After = rb.After
	r.DateHeaders = rb.DateHeaders
	r.selectorBlock = sb
	r.Remove = rb.Remove
	return nil
}

const (
	searchMethodPost = "post"
	searchMethodGet  = "get"
)

type searchBlock struct {
	Path   string          `yaml:"path"`
	Method string          `yaml:"method"`
	Inputs inputsBlock     `yaml:"inputs,omitempty"`
	Rows   rowsBlock       `yaml:"rows"`
	Fields fieldsListBlock `yaml:"fields"`
}

type capabilitiesBlock struct {
	CategoryMap categoryMap
	SearchModes []torznab.SearchMode
}

// UnmarshalYAML implements the Unmarshaller interface.
func (c *capabilitiesBlock) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var intermediate struct {
		Categories map[string]string        `yaml:"categories"`
		Modes      map[string]stringorslice `yaml:"modes"`
	}

	if err := unmarshal(&intermediate); err == nil {
		c.CategoryMap = categoryMap{}

		for id, catName := range intermediate.Categories {
			matchedCat := false
			for _, cat := range torznab.AllCategories {
				if cat.Name == catName {
					c.CategoryMap[id] = cat
					matchedCat = true
					break
				}
			}
			if !matchedCat {
				return fmt.Errorf("Unknown category %q", catName)
			}
		}

		c.SearchModes = []torznab.SearchMode{}

		for key, supported := range intermediate.Modes {
			c.SearchModes = append(c.SearchModes, torznab.SearchMode{key, true, supported})
		}

		return nil
	}

	return errors.New("Failed to unmarshal capabilities block")
}

func (c *capabilitiesBlock) ToTorznab() torznab.Capabilities {
	caps := torznab.Capabilities{
		Categories:  c.CategoryMap.Categories(),
		SearchModes: []torznab.SearchMode{},
	}

	// All indexers support search
	caps.SearchModes = append(caps.SearchModes, torznab.SearchMode{
		Key:             "search",
		Available:       true,
		SupportedParams: []string{"q"},
	})

	// Some support TV
	if caps.HasTVShows() {
		caps.SearchModes = append(caps.SearchModes, torznab.SearchMode{
			Key:             "tv-search",
			Available:       true,
			SupportedParams: []string{"q", "season", "ep"},
		})
	}

	// Some support Movies
	if caps.HasMovies() {
		caps.SearchModes = append(caps.SearchModes, torznab.SearchMode{
			Key:             "movie-search",
			Available:       true,
			SupportedParams: []string{"q"},
		})
	}

	return caps
}

type stringorslice []string

func (s *stringorslice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var stringType string
	if err := unmarshal(&stringType); err == nil {
		*s = stringorslice{stringType}
		return nil
	}

	var sliceType []string
	if err := unmarshal(&sliceType); err == nil {
		*s = stringorslice(sliceType)
		return nil
	}

	return errors.New("Failed to unmarshal stringorslice")
}

type ratioBlock struct {
	selectorBlock
	Path string `yaml:"path"`
}

func (r *ratioBlock) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var sb selectorBlock
	if err := unmarshal(&sb); err != nil {
		return errors.New("Failed to unmarshal ratioBlock")
	}

	var rb struct {
		Path string `yaml:"path"`
	}
	if err := unmarshal(&rb); err != nil {
		return errors.New("Failed to unmarshal ratioBlock")
	}

	r.selectorBlock = sb
	r.Path = rb.Path
	return nil
}
