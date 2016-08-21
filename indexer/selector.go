package indexer

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

type filterBlock struct {
	Name string      `yaml:"name"`
	Args interface{} `yaml:"args"`
}

type selectorBlock struct {
	Selector  string        `yaml:"selector"`
	Attribute string        `yaml:"attribute,omitempty"`
	Remove    string        `yaml:"remove,omitempty"`
	Filters   []filterBlock `yaml:"filters,omitempty"`
}

func (s *selectorBlock) Match(selection *goquery.Selection) bool {
	return selection.Find(s.Selector).Length() > 0
}

func (s *selectorBlock) Text(selection *goquery.Selection) (string, error) {
	result := selection.Find(s.Selector)
	if result.Length() == 0 {
		return "", nil
	}

	html, _ := result.Html()
	log.Printf("Selector %q matched %q", s.Selector, html)

	if s.Remove != "" {
		result.Find(s.Remove).Remove()
	}

	output := result.Text()

	if s.Attribute != "" {
		val, exists := result.Attr(s.Attribute)
		if !exists {
			return "", fmt.Errorf("Requested attribute %q doesn't exist", s.Attribute)
		}
		output = val
	}

	for _, f := range s.Filters {
		log.Printf("Applying filter %s(%#v) to %q", f.Name, f.Args, output)

		var err error
		output, err = dispatchFilter(f.Name, f.Args, output)
		if err != nil {
			return "", err
		}
	}

	log.Printf("Final text is %q", output)
	return output, nil
}

func (s *selectorBlock) IsEmpty() bool {
	return s.Selector == ""
}
