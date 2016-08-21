package indexer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/cardigann/cardigann/torznab"
	"github.com/headzoo/surf/browser"
	"github.com/shibukawa/configdir"

	"gopkg.in/yaml.v2"
)

type IndexerDefinition struct {
	Site         string            `yaml:"site"`
	Name         string            `yaml:"name"`
	Description  string            `yaml:"description"`
	Language     string            `yaml:"language"`
	Links        Stringorslice     `yaml:"links"`
	Capabilities CapabilitiesBlock `yaml:"caps"`
	Login        LoginBlock        `yaml:"login"`
	Search       SearchBlock       `yaml:"search"`
}

func (i *IndexerDefinition) applyDefaults() {
	if i.Language == "" {
		i.Language = "en-us"
	}

	if i.Login.FormSelector == "" {
		i.Login.FormSelector = "form"
	}
}

// Stringorslice represents a string or an array of strings.
type Stringorslice []string

// UnmarshalYAML implements the Unmarshaller interface.
func (s *Stringorslice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var stringType string
	if err := unmarshal(&stringType); err == nil {
		*s = Stringorslice{stringType}
		return nil
	}

	var sliceType []string
	if err := unmarshal(&sliceType); err == nil {
		*s = Stringorslice(sliceType)
		return nil
	}

	return errors.New("Failed to unmarshal Stringorslice")
}

type InputsBlock map[string]string

type SelectorBlock struct {
	Selector  string   `yaml:"selector"`
	Attribute string   `yaml:"attribute"`
	Filters   []Filter `yaml:"filters,omitempty"`
}

type Filter struct {
	Name string      `yaml:"name"`
	Args interface{} `yaml:"args"`
}

func (s *SelectorBlock) Match(selection *goquery.Selection) (string, bool) {
	result := selection.Find(s.Selector)

	if result.Length() == 0 {
		return "", false
	}

	var output string

	if s.Attribute != "" {
		val, exists := result.Attr(s.Attribute)
		if !exists {
			return "", false
		}
		output = val
	} else {
		output = result.Text()
	}

	for _, f := range s.Filters {
		log.Printf("Applying filter %s(%#v)", f.Name, f.Args)

		var err error
		output, err = dispatchFilter(f.Name, f.Args, output)
		if err != nil {
			panic(err)
		}
	}

	return output, true
}

func (s *SelectorBlock) IsZero() bool {
	return s.Selector == ""
}

type ErrorBlock struct {
	Path     string        `yaml:"path"`
	Selector string        `yaml:"selector"`
	Message  SelectorBlock `yaml:"message"`
}

func (e *ErrorBlock) Match(browser browser.Browsable) (string, bool) {
	if e.Path != "" {
		if e.Path != browser.Url().Path {
			return "", false
		}
	}

	if e.Selector != "" {
		result := browser.Find(e.Selector)

		if result.Length() == 0 {
			return "", false
		}

		if e.Message.IsZero() {
			return result.Text(), true
		}
	}

	return e.Message.Match(browser.Dom())
}

type LoginBlock struct {
	Path         string      `yaml:"path"`
	FormSelector string      `yaml:"form"`
	Inputs       InputsBlock `yaml:"inputs,omitempty"`
	Error        ErrorBlock  `yaml:"error,omitempty"`
}

type FieldsBlock map[string]SelectorBlock

type SearchBlock struct {
	Path   string        `yaml:"path"`
	Inputs InputsBlock   `yaml:"inputs,omitempty"`
	Rows   SelectorBlock `yaml:"rows"`
	Fields FieldsBlock   `yaml:"fields"`
}

type CapabilitiesBlock torznab.Capabilities

// UnmarshalYAML implements the Unmarshaller interface.
func (c *CapabilitiesBlock) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var intermediate struct {
		Categories map[int]string           `yaml:"categories"`
		Modes      map[string]Stringorslice `yaml:"modes"`
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

func ParseDefinitionFile(f *os.File) (*IndexerDefinition, error) {
	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, err
	}

	return ParseDefinition(b)
}

func ParseDefinition(src []byte) (*IndexerDefinition, error) {
	def := IndexerDefinition{
		Capabilities: CapabilitiesBlock{},
		Login: LoginBlock{
			Inputs: InputsBlock{},
		},
		Search: SearchBlock{},
	}

	if err := yaml.Unmarshal(src, &def); err != nil {
		return nil, err
	}

	def.applyDefaults()
	return &def, nil
}

func LoadDefinition(key string) (*IndexerDefinition, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cd := configdir.New("cardigann", "cardigann")
	cd.LocalPath = cwd

	fileName := key + ".yml"
	folder := cd.QueryFolderContainsFile(fileName)
	if folder == nil {
		return nil, errors.New("Failed to find " + fileName)
	}

	data, err := folder.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return ParseDefinition(data)
}
