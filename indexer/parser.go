package indexer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cardigann/cardigann/torznab"
	"github.com/headzoo/surf/browser"

	"gopkg.in/yaml.v2"
)

type IndexerDefinition struct {
	Site         string            `yaml:"site"`
	Name         string            `yaml:"name"`
	Description  string            `yaml:"description"`
	Language     string            `yaml:"language"`
	Links        stringorslice     `yaml:"links"`
	Capabilities capabilitiesBlock `yaml:"caps"`
	Login        loginBlock        `yaml:"login"`
	Search       searchBlock       `yaml:"search"`
}

func ParseDefinitionFile(f *os.File) (*IndexerDefinition, error) {
	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, err
	}

	return ParseDefinition(b)
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

	if def.Login.Error.Message.IsEmpty() && def.Login.Error.Selector != "" {
		def.Login.Error.Message.Selector = def.Login.Error.Selector
	}

	return &def, nil
}

type inputsBlock map[string]string

type errorBlock struct {
	Path     string        `yaml:"path"`
	Selector string        `yaml:"selector"`
	Message  selectorBlock `yaml:"message"`
}

func (e *errorBlock) MatchPage(browser browser.Browsable) bool {
	if e.Path != "" {
		if e.Path != browser.Url().Path {
			return false
		}
	}

	if e.Selector != "" {
		return browser.Find(e.Selector).Length() > 0
	}

	return false
}

type loginBlock struct {
	Path         string      `yaml:"path"`
	FormSelector string      `yaml:"form"`
	Inputs       inputsBlock `yaml:"inputs,omitempty"`
	Error        errorBlock  `yaml:"error,omitempty"`
}

type fieldsBlock map[string]selectorBlock

type searchBlock struct {
	Path   string        `yaml:"path"`
	Inputs inputsBlock   `yaml:"inputs,omitempty"`
	Rows   selectorBlock `yaml:"rows"`
	Fields fieldsBlock   `yaml:"fields"`
}

type capabilitiesBlock torznab.Capabilities

// UnmarshalYAML implements the Unmarshaller interface.
func (c *capabilitiesBlock) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var intermediate struct {
		Categories map[int]string           `yaml:"categories"`
		Modes      map[string]stringorslice `yaml:"modes"`
	}

	if err := unmarshal(&intermediate); err == nil {
		c.Categories = torznab.CategoryMapping{}

		for id, catName := range intermediate.Categories {
			matchedCat := false
			for _, cat := range torznab.AllCategories {
				if cat.Name == catName {
					c.Categories[id] = cat
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

	return errors.New("Failed to unmarshal CapabilitiesBlock")
}

// Stringorslice represents a string or an array of strings.
type stringorslice []string

// UnmarshalYAML implements the Unmarshaller interface.
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
