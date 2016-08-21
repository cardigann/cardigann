package indexer

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"time"
)

func dispatchFilter(name string, args interface{}, value string) (string, error) {
	switch name {
	case "querystring":
		param, ok := args.(string)
		if !ok {
			return "", fmt.Errorf("Filter %q requires a string argument", name)
		}
		return filterQueryString(param, value)

	case "dateparse":
		format, ok := args.(string)
		if !ok {
			return "", fmt.Errorf("Filter %q requires a string argument", name)
		}
		return filterDateParse(format, value)

	case "regexp":
		pattern, ok := args.(string)
		if !ok {
			return "", fmt.Errorf("Filter %q requires a string argument", name)
		}
		return filterRegexp(pattern, value)
	}

	return "", errors.New("Unknown filter " + name)
}

func filterQueryString(param string, value string) (string, error) {
	u, err := url.Parse(value)
	if err != nil {
		return "", err
	}
	return u.Query().Get(param), nil
}

func filterDateParse(format string, value string) (string, error) {
	t, err := time.Parse(format, value)
	if err != nil {
		return "", err
	}
	return t.Format(time.RFC1123Z), nil
}

func filterRegexp(pattern string, value string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	matches := re.FindStringSubmatch(value)

	if len(matches) == 0 {
		return "", errors.New("No matches found for pattern")
	}

	if len(matches) > 1 {
		return matches[1], nil
	}

	return matches[0], nil
}
