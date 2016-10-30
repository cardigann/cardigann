package indexer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
	"github.com/yosssi/gohtml"
)

type filterBlock struct {
	Name string      `yaml:"name"`
	Args interface{} `yaml:"args"`
}

type selectorBlock struct {
	Selector  string            `yaml:"selector"`
	TextVal   string            `yaml:"text"`
	Attribute string            `yaml:"attribute,omitempty"`
	Remove    string            `yaml:"remove,omitempty"`
	Filters   []filterBlock     `yaml:"filters,omitempty"`
	Case      map[string]string `yaml:"case,omitempty"`
}

func (s *selectorBlock) Match(selection *goquery.Selection) bool {
	return !s.IsEmpty() && (selection.Find(s.Selector).Length() > 0 || s.TextVal != "")
}

func (s *selectorBlock) MatchText(from *goquery.Selection) (string, error) {
	if s.TextVal != "" {
		return s.TextVal, nil
	}
	if s.Selector != "" {
		result := from.Find(s.Selector)
		if result.Length() == 0 {
			return "", fmt.Errorf("Failed to match selector %q", s.Selector)
		}
		return s.Text(result)
	}
	return s.Text(from)
}

func (s *selectorBlock) Text(el *goquery.Selection) (string, error) {
	if s.TextVal != "" {
		return s.applyFilters(s.TextVal)
	}

	if s.Remove != "" {
		el.Find(s.Remove).Remove()
	}

	if s.Case != nil {
		filterLogger.
			WithFields(logrus.Fields{"case": s.Case}).
			Debugf("Applying case to selection")
		for pattern, value := range s.Case {
			if el.Is(pattern) || el.Has(pattern).Length() >= 1 {
				return s.applyFilters(value)
			}
		}
		return "", errors.New("None of the cases match")
	}

	html, _ := goquery.OuterHtml(el)
	filterLogger.
		WithFields(logrus.Fields{"html": gohtml.Format(html)}).
		Debugf("Extracting text from selection")

	output := strings.TrimSpace(el.Text())

	if s.Attribute != "" {
		val, exists := el.Attr(s.Attribute)
		if !exists {
			return "", fmt.Errorf("Requested attribute %q doesn't exist", s.Attribute)
		}
		output = val
	}

	return s.applyFilters(output)
}

func (s *selectorBlock) applyFilters(val string) (string, error) {
	for _, f := range s.Filters {
		filterLogger.
			WithFields(logrus.Fields{"args": f.Args, "before": val}).
			Debugf("Applying filter %s", f.Name)

		var err error
		val, err = invokeFilter(f.Name, f.Args, val)
		if err != nil {
			return "", err
		}
	}

	return val, nil
}

func (s *selectorBlock) IsEmpty() bool {
	return s.Selector == "" && s.TextVal == ""
}

func (s *selectorBlock) String() string {
	switch {
	case s.Selector != "":
		return fmt.Sprintf("Selector(%s)", s.Selector)
	case s.TextVal != "":
		return fmt.Sprintf("Text(%s)", s.TextVal)
	default:
		return "Empty"
	}
}
