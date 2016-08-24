package indexer

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"
	"time"
)

const (
	filterTimeFormat = time.RFC1123Z
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
	return t.Format(filterTimeFormat), nil
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

func filterTimeAgo(src string, now time.Time) (string, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	var tok rune
	for tok != scanner.EOF {
		tok = s.Scan()
		switch s.TokenText() {
		case ",", "ago", "", "and":
			continue
		}

		f, err := strconv.ParseFloat(s.TokenText(), 64)
		if err != nil {
			return "", fmt.Errorf(
				"Failed to parse duration %q in time format at %s", s.TokenText(), s.Pos())
		}

		tok = s.Scan()
		if tok == scanner.EOF {
			return "", fmt.Errorf(
				"Expected a time unit at %s", s.TokenText(), s.Pos())
		}

		switch strings.TrimSuffix(s.TokenText(), "s") {
		case "year":
			now = now.AddDate(-int(f), 0, 0)
		case "month":
			now = now.AddDate(0, -int(f), 0)
		case "week":
			now = now.Add(-time.Hour * 24 * 7)
		case "day":
			now = now.AddDate(0, 0, -int(f))
		case "hour":
			now = now.Add(-time.Hour)
		case "minute":
			now = now.Add(-time.Minute)
		default:
			return "", fmt.Errorf("Unsupporting unit of time %q", s.TokenText())
		}
	}

	return now.Format(filterTimeFormat), nil
}
