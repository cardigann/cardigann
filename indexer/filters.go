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
	"unicode"

	"github.com/Sirupsen/logrus"
	"github.com/cardigann/cardigann/logger"
)

const (
	filterTimeFormat = time.RFC1123Z
)

var (
	filterLogger = logger.Logger
)

func invokeFilter(name string, args interface{}, value string) (string, error) {
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

	case "split":
		sep, ok := (args.([]interface{}))[0].(string)
		if !ok {
			return "", fmt.Errorf("Filter %q requires a string argument at idx 0", name)
		}
		pos, ok := (args.([]interface{}))[1].(int)
		if !ok {
			return "", fmt.Errorf("Filter %q requires an int argument at idx 1", name)
		}
		return filterSplit(sep, pos, value)

	case "replace":
		from, ok := (args.([]interface{}))[0].(string)
		if !ok {
			return "", fmt.Errorf("Filter %q requires a string argument at idx 0", name)
		}
		to, ok := (args.([]interface{}))[1].(string)
		if !ok {
			return "", fmt.Errorf("Filter %q requires a string argument at idx 1", name)
		}
		return strings.Replace(value, from, to, -1), nil

	case "trim":
		cutset, ok := args.(string)
		if !ok {
			return "", fmt.Errorf("Filter %q requires a string argument at idx 0", name)
		}
		return strings.Trim(value, cutset), nil

	case "timeago":
		return filterTimeAgo(value, time.Now())

	case "reltime":
		format, ok := args.(string)
		if !ok {
			return "", fmt.Errorf("Filter %q requires a string argument", name)
		}
		return filterRelTime(value, format, time.Now())
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

func filterSplit(sep string, pos int, value string) (string, error) {
	frags := strings.Split(value, sep)
	if pos < 0 {
		pos = len(frags) + pos
	}
	return frags[pos], nil
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

	filterLogger.WithFields(logrus.Fields{"matches": matches}).Debug("Regex matched")

	if len(matches) > 1 {
		return matches[1], nil
	}

	return matches[0], nil
}

func splitDecimalStr(s string) (int, float64, error) {
	if parts := strings.SplitN(s, ".", 2); len(parts) == 2 {
		i, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, err
		}
		f, err := strconv.ParseFloat("0."+parts[1], 64)
		if err != nil {
			return 0, 0, err
		}
		return i, f, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, 0, err
	}
	return i, 0, nil
}

func normalizeSpace(s string) string {
	return strings.TrimSpace(strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		return r
	}, s))
}

func filterTimeAgo(src string, now time.Time) (string, error) {
	normalized := normalizeSpace(src)
	normalized = strings.ToLower(normalized)

	var s scanner.Scanner
	s.Init(strings.NewReader(normalized))
	var tok rune
	for tok != scanner.EOF {
		tok = s.Scan()

		switch s.TokenText() {
		case ",", "ago", "", "and":
			continue
		}

		v, fraction, err := splitDecimalStr(s.TokenText())
		if err != nil {
			return "", fmt.Errorf(
				"Failed to parse decimal time %q in time format at %s", s.TokenText(), s.Pos())
		}

		tok = s.Scan()
		if tok == scanner.EOF {
			return "", fmt.Errorf(
				"Expected a time unit at %s", s.TokenText(), s.Pos())
		}

		switch strings.TrimSuffix(s.TokenText(), "s") {
		case "year":
			now = now.AddDate(-v, 0, 0)
			if fraction > 0 {
				now = now.Add(time.Duration(float64(now.AddDate(-1, 0, 0).Sub(now)) * fraction))
			}
		case "month":
			now = now.AddDate(0, -v, 0)
			if fraction > 0 {
				now = now.Add(time.Duration(float64(now.AddDate(0, -1, 0).Sub(now)) * fraction))
			}
		case "week":
			now = now.AddDate(0, 0, -7)
			if fraction > 0 {
				now = now.Add(time.Duration(float64(now.AddDate(0, 0, -7).Sub(now)) * fraction))
			}
		case "day":
			now = now.AddDate(0, 0, -v)
			if fraction > 0 {
				now = now.Add(time.Minute * -time.Duration(fraction*1440))
			}
		case "hour":
			now = now.Add(time.Hour * -time.Duration(v))
			if fraction > 0 {
				now = now.Add(time.Second * -time.Duration(fraction*3600))
			}
		case "minute":
			now = now.Add(time.Minute * -time.Duration(v))
			if fraction > 0 {
				now = now.Add(time.Second * -time.Duration(fraction*60))
			}
		case "second":
			now = now.Add(time.Second * -time.Duration(v))
		default:
			return "", fmt.Errorf("Unsupporting unit of time %q", s.TokenText())
		}
	}

	return now.Format(filterTimeFormat), nil
}

func filterRelTime(src string, format string, now time.Time) (string, error) {
	out := src

	for from, to := range map[string]string{
		"today":     now.Format(format),
		"Today":     now.Format(format),
		"yesterday": now.AddDate(0, 0, -1).Format(format),
		"Yesterday": now.AddDate(0, 0, -1).Format(format),
	} {
		out = strings.Replace(out, from, to, -1)
	}

	return out, nil
}
