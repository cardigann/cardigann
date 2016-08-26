package indexer

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
)

type filterBlock struct {
	Name string      `yaml:"name"`
	Args interface{} `yaml:"args"`
}

type selectorBlock struct {
	Selector  string        `yaml:"selector"`
	TextVal   string        `yaml:"text"`
	Attribute string        `yaml:"attribute,omitempty"`
	Remove    string        `yaml:"remove,omitempty"`
	Filters   []filterBlock `yaml:"filters,omitempty"`
}

func (s *selectorBlock) Match(selection *goquery.Selection) bool {
	return selection.Find(s.Selector).Length() > 0
}

func (s *selectorBlock) Text(selection *goquery.Selection) (string, error) {
	output := s.TextVal

	if s.Selector != "" {
		result := selection.Find(s.Selector)
		if result.Length() == 0 {
			return "", nil
		}

		html, _ := result.Html()
		filterLogger.
			WithFields(logrus.Fields{"selector": s.Selector, "html": strings.TrimSpace(html)}).
			Debugf("Selector matched %d elements", result.Length())

		if s.Remove != "" {
			result.Find(s.Remove).Remove()
		}

		output = strings.TrimSpace(result.Text())

		if s.Attribute != "" {
			val, exists := result.Attr(s.Attribute)
			if !exists {
				return "", fmt.Errorf("Requested attribute %q doesn't exist", s.Attribute)
			}
			output = val
		}
	}

	for _, f := range s.Filters {
		filterLogger.
			WithFields(logrus.Fields{"args": f.Args, "before": output}).
			Debugf("Applying filter %s", f.Name)

		var err error
		output, err = invokeFilter(f.Name, f.Args, output)
		if err != nil {
			return "", err
		}
	}

	return output, nil
}

func (s *selectorBlock) IsEmpty() bool {
	return s.Selector == ""
}
