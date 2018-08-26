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
	"github.com/bcampbell/fuzzytime"
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

	case "timeparse", "dateparse":
		if args == nil {
			return filterDateParse(nil, value)
		}
		if layout, ok := args.(string); ok {
			return filterDateParse([]string{layout}, value)
		}
		return "", fmt.Errorf("Filter argument type %T was invalid", args)

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

	case "append":
		str, ok := args.(string)
		if !ok {
			return "", fmt.Errorf("Filter %q requires a string argument at idx 0", name)
		}
		return value + str, nil

	case "prepend":
		str, ok := args.(string)
		if !ok {
			return "", fmt.Errorf("Filter %q requires a string argument at idx 0", name)
		}
		return str + value, nil

	case "timeago", "fuzzytime", "reltime":
		return filterFuzzyTime(value, time.Now())
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

func filterDateParse(layouts []string, value string) (string, error) {
	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t.Format(filterTimeFormat), nil
		}
	}
	return "", fmt.Errorf("No matching date pattern for %s", value)
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

var (
	timeAgoRegexp     = regexp.MustCompile(`(?i)\bago`)
	todayRegexp       = regexp.MustCompile(`(?i)\btoday([\s,]+|$)`)
	tomorrowRegexp    = regexp.MustCompile(`(?i)\btomorrow([\s,]+|$)`)
	yesterdayRegexp   = regexp.MustCompile(`(?i)\byesterday([\s,]+|$)`)
	missingYearRegexp = regexp.MustCompile(`^\d{1,2}-\d{1,2}\b`)
)

func normalizeNumber(s string) string {
	normalized := normalizeSpace(s)
	normalized = strings.Trim(s, "-")
	normalized = strings.Replace(s, ",", "", -1)

	if normalized == "" {
		normalized = "0"
	}

	return normalized
}

func normalizeSpace(s string) string {
	return strings.TrimSpace(strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		return r
	}, s))
}

func parseTimeAgo(src string, now time.Time) (time.Time, error) {
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
			return now, fmt.Errorf(
				"failed to parse decimal time %q in time format at %s", s.TokenText(), s.Pos())
		}

		tok = s.Scan()
		if tok == scanner.EOF {
			return now, fmt.Errorf(
				"expected a time unit at %s:%v", s.TokenText(), s.Pos())
		}

		unit := s.TokenText()
		if unit != "s" {
			unit = strings.TrimSuffix(s.TokenText(), "s")
		}

		switch unit {
		case "year", "yr", "y":
			now = now.AddDate(-v, 0, 0)
			if fraction > 0 {
				now = now.Add(time.Duration(float64(now.AddDate(-1, 0, 0).Sub(now)) * fraction))
			}
		case "month", "mnth", "mo":
			now = now.AddDate(0, -v, 0)
			if fraction > 0 {
				now = now.Add(time.Duration(float64(now.AddDate(0, -1, 0).Sub(now)) * fraction))
			}
		case "week", "wk", "w":
			now = now.AddDate(0, 0, -7)
			if fraction > 0 {
				now = now.Add(time.Duration(float64(now.AddDate(0, 0, -7).Sub(now)) * fraction))
			}
		case "day", "d":
			now = now.AddDate(0, 0, -v)
			if fraction > 0 {
				now = now.Add(time.Minute * -time.Duration(fraction*1440))
			}
		case "hour", "hr", "h":
			now = now.Add(time.Hour * -time.Duration(v))
			if fraction > 0 {
				now = now.Add(time.Second * -time.Duration(fraction*3600))
			}
		case "minute", "min", "m":
			now = now.Add(time.Minute * -time.Duration(v))
			if fraction > 0 {
				now = now.Add(time.Second * -time.Duration(fraction*60))
			}
		case "second", "sec", "s":
			now = now.Add(time.Second * -time.Duration(v))
		default:
			return now, fmt.Errorf("Unsupporting unit of time %q", unit)
		}
	}

	return now, nil
}

func parseFuzzyTime(src string, now time.Time) (time.Time, error) {
	if timeAgoRegexp.MatchString(src) {
		t, err := parseTimeAgo(src, now)
		if err != nil {
			return t, fmt.Errorf("error parsing time ago %q: %v", src, err)
		}
		return t, nil
	}

	normalized := normalizeSpace(src)

	out := todayRegexp.ReplaceAllLiteralString(normalized, now.Format("Mon, 02 Jan 2006 "))
	out = tomorrowRegexp.ReplaceAllLiteralString(out, now.AddDate(0, 0, 1).Format("Mon, 02 Jan 2006 "))
	out = yesterdayRegexp.ReplaceAllLiteralString(out, now.AddDate(0, 0, -1).Format("Mon, 02 Jan 2006 "))

	if m := missingYearRegexp.FindStringSubmatch(out); len(m) > 0 {
		out = missingYearRegexp.ReplaceAllLiteralString(src, m[0]+now.Format("-2006"))
	}

	dt, _, err := fuzzytime.USContext.Extract(out)
	if err != nil {
		return time.Time{}, fmt.Errorf("error extracting date from %q: %v", out, err)
	}

	if dt.Time.Empty() {
		dt.Time.SetHour(0)
		dt.Time.SetMinute(0)
	}

	if !dt.HasFullDate() {
		return time.Time{}, fmt.Errorf("found only partial date %v", dt.ISOFormat())
	}

	if !dt.Time.HasSecond() {
		dt.Time.SetSecond(0)
	}

	if !dt.HasTZOffset() {
		dt.Time.SetTZOffset(0)
	}

	return time.Parse("2006-01-02T15:04:05Z07:00", dt.ISOFormat())
}

func filterFuzzyTime(src string, now time.Time) (string, error) {
	t, err := parseFuzzyTime(src, now)
	if err != nil {
		return "", fmt.Errorf("error parsing fuzzy time %q: %v", src, err)
	}
	return t.Format(filterTimeFormat), nil
}
