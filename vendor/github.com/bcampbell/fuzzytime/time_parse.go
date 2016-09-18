package fuzzytime

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// match one of:
//  named timezone (eg, BST, NZDT etc)
//  Z
//  +hh:mm, +hhmm, or +hh
//  -hh:mm, -hhmm, or -hh
var tzPat = `(?i)(?P<tz>Z|[A-Z]{2,5}|(([-+])(\d{2})((:?)(\d{2}))?))`
var ampmPat = `(?i)(?:(?P<am>(am|a[.]m[.]))|(?P<pm>(pm|p[.]m[.])))`

var timeCrackers = []*regexp.Regexp{
	// "4:48PM GMT"
	regexp.MustCompile(`(?i)(?P<hour>\d{1,2})[:.](?P<min>\d{2})(?:[:.](?P<sec>\d{2}))?[\s\p{Z}]*` + ampmPat + `[\s\p{Z}]*` + tzPat),

	// "3:34PM"
	// "10:42 am"
	regexp.MustCompile(`(?i)\b(?P<hour>\d{1,2})[:.](?P<min>\d{2})(?:[:.](?P<sec>\d{2}))?[\s\p{Z}]*` + ampmPat),

	// "13:21:36 GMT"
	// "15:29 GMT"
	// "12:35:44+00:00"
	// "23:59:59.9942+01:00"
	regexp.MustCompile(`(?i)(?:\b|T)(?P<hour>\d{1,2})[:](?P<min>\d{2})(?:[:](?P<sec>\d{2})(?P<fractional>[.]\d+)?)?[\s\p{Z}]*` + tzPat),

	// "00.01 BST"
	regexp.MustCompile(`(?i)(?:\b|T)(?P<hour>\d{1,2})[.](?P<min>\d{2})(?:[.](?P<sec>\d{2}))?[\s\p{Z}]*` + tzPat),

	// "14:21:01"
	// "14:21"
	// "23:59:59.9942"
	regexp.MustCompile(`(?i)(?:\b|T)(?P<hour>\d{1,2})[:](?P<min>\d{2})(?:[:](?P<sec>\d{2})(?P<fractional>[.]\d+)?)?(?:[^\d]|\z)`),
}

// ExtractTime tries to parse a time from a string.
// It returns a Time and a Span indicating which part of string matched.
// Time and Span may be empty, indicating no time was found.
// An error will be returned if a time is found but cannot be correctly parsed.
// If error is not nil time the returned time and span will both be empty
func (ctx *Context) ExtractTime(s string) (Time, Span, error) {
	for _, pat := range timeCrackers {
		names := pat.SubexpNames()
		matchSpans := pat.FindStringSubmatchIndex(s)
		if matchSpans == nil {
			continue
		}

		var hour, minute, second int = -1, -1, -1
		var am, pm bool = false, false

		var gotTZ = false
		var tzOffset int
		var fail, err error
		for i, name := range names {
			start, end := matchSpans[i*2], matchSpans[(i*2)+1]
			if start == end {
				continue
			}
			var sub string
			if start >= 0 && end >= 0 {
				sub = strings.ToLower(s[start:end])
			}

			switch name {
			case "hour":
				hour, err = strconv.Atoi(sub)
				if err != nil {
					fail = err
					break
				}
				if hour < 0 || hour > 23 {
					fail = errors.New("bad hour value")
					break
				}

			case "min":
				minute, err = strconv.Atoi(sub)
				if err != nil {
					fail = err
					break
				}
				if minute < 0 || minute > 59 {
					fail = errors.New("bad minute value")
					break
				}
			case "sec":
				second, err = strconv.Atoi(sub)
				if err != nil {
					fail = err
					break
				}
				if second < 0 || second > 59 {
					fail = errors.New("bad seconds value")
					break
				}
			case "am":
				am = true
			case "pm":
				pm = true
			case "tz":
				offset, err := ctx.parseTZ(sub)
				if err != nil {
					break
					//return Time{}, Span{}, err
				}
				tzOffset = offset
				gotTZ = true
			case "fractional":
				// ignore fractional seconds for now
			}

		}

		if fail != nil {
			break
		}

		// got enough to accept?
		if hour >= 0 && minute >= 0 {
			if pm && (hour >= 1) && (hour <= 11) {
				hour += 12
			}
			if am && (hour == 12) {
				hour -= 12
			}

			ft := Time{}
			ft.SetHour(hour)
			ft.SetMinute(minute)
			if second >= 0 {
				ft.SetSecond(second)
			}
			if gotTZ {
				ft.SetTZOffset(tzOffset)
			}
			var span = Span{matchSpans[0], matchSpans[1]}
			return ft, span, nil
		}
	}

	// nothing. Just return an empty time and span
	return Time{}, Span{}, nil
}

func (ctx *Context) parseTZ(s string) (int, error) {
	s = strings.ToUpper(s)
	// try as an ISO 8601-style offset ("+01:30" etc)
	offset, err := TZToOffset(s)
	if err != nil {
		// nope, try resolving as a named timezone via the context
		offset, err = ctx.TZResolver(s)
	}
	return offset, err
}
