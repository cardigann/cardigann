package fuzzytime

import (
	"regexp"
	"strconv"
	"strings"
)

// dateCracker is a set of regexps for various date formats
// order is important(ish) - want to match as much of the string as we can
var dateCrackers = []*regexp.Regexp{
	//"Tuesday 16 December 2008"
	//"Tue 29 Jan 08"
	//"Monday, 22 October 2007"
	//"Tuesday, 21st January, 2003"
	regexp.MustCompile(`(?i)(?P<dayname>\w{3,})[.,\s\p{Z}]+(?P<day>\d{1,2})(?:st|nd|rd|th)?[\s\p{Z}]+(?P<month>\w{3,})[.,\s\p{Z}]+(?P<year>(\d{4})|(\d{2}))`),

	// "Friday    August    11, 2006"
	// "Tuesday October 14 2008"
	// "Thursday August 21 2008"
	// "Monday, May. 17, 2010"
	regexp.MustCompile(`(?i)(?P<dayname>\w{3,})[.,\s\p{Z}]+(?P<month>\w{3,})[.,\s\p{Z}]+(?P<day>\d{1,2})(?:st|nd|rd|th)?[.,\s\p{Z}]+(?P<year>(\d{4})|(\d{2}))`),

	// "9 Sep 2009", "09 Sep, 2009", "01 May 10"
	// "23rd November 2007", "22nd May 2008"
	regexp.MustCompile(`(?i)(?P<day>\d{1,2})(?:st|nd|rd|th)?[\s\p{Z}]+(?P<month>\w{3,})[.,\s\p{Z}]+(?P<year>(\d{4})|(\d{2}))`),

	// "Mar 3, 2007", "Jul 21, 08", "May 25 2010", "May 25th 2010", "February 10 2008"
	regexp.MustCompile(`(?i)(?P<month>\w{3,})[.,\s\p{Z}]+(?P<day>\d{1,2})(?:st|nd|rd|th)?[.,\s\p{Z}]+(?P<year>(\d{4})|(\d{2}))`),

	// "2010-04-02"
	regexp.MustCompile(`(?i)(?P<year>\d{4})-(?P<month>\d{1,2})-(?P<day>\d{1,2})`),

	// "2007/03/18"
	regexp.MustCompile(`(?i)(?P<year>\d{4})/(?P<month>\d{1,2})/(?P<day>\d{1,2})`),

	// "09-Apr-2007", "09-Apr-07"
	regexp.MustCompile(`(?i)(?P<day>\d{1,2})-(?P<month>\w{3,})-(?P<year>(\d{4})|(\d{2}))`),

	// "May 2011"
	regexp.MustCompile(`(?i)(?P<month>\w{3,})[\s\p{Z}]+(?P<year>\d{4})`),

	// ambiguous formats
	// "11/02/2008"
	// "11-02-2008"
	// "11.02.2008"
	regexp.MustCompile(`(?i)(?P<x1>\d{1,2})[/.-](?P<x2>\d{1,2})[/.-](?P<year>\d{4})`),
	// even more ambiguous
	// eg:  japan uses yy/mm/dd
	// 11/2/10
	// 11-02-10
	// 11.02.10
	regexp.MustCompile(`(?i)(?P<x1>\d{1,2})[/-](?P<x2>\d{1,2})[/-](?P<x3>\d{2})`),
	/*
	   # TODO:
	   # year/month only

	   # "May/June 2011" (common for publications) - just use second month
	   r'(?P<cruftmonth>\w{3,})/(?P<month>\w{3,})[\s\p{Z}]+(?P<year>\d{4})',
	*/

	// Missing year, eg
	// Thu April 24th
	regexp.MustCompile(`(?i)(?P<dayname>\w{3,})[.,\s\p{Z}]+(?P<month>\w{3,})[.,\s\p{Z}]+(?P<day>\d{1,2})(?:st|nd|rd|th)?`),

	// April 24th
	regexp.MustCompile(`(?i)(?P<month>\w{3,})[.,\s\p{Z}]+(?P<day>\d{1,2})(?:st|nd|rd|th)?`),
}

// ExtendYear extends 2-digit years into 4 digits.
// the rules used:
// 00-69 => 2000-2069
// 70-99 => 1970-1999
func ExtendYear(year int) int {
	if year < 70 {
		return 2000 + year
	}
	if year < 100 {
		return 1900 + year
	}
	return year
}

// ExtractDate tries to parse a date from a string.
// It returns a Date and Span indicating which part of string matched.
// If an error occurs, an empty Date will be returned.
func (ctx *Context) ExtractDate(s string) (Date, Span, error) {

	for _, pat := range dateCrackers {
		fd := Date{}
		span := Span{}
		names := pat.SubexpNames()
		matchSpans := pat.FindStringSubmatchIndex(s)
		if matchSpans == nil {
			continue
		}

		var fail bool

		unknowns := make([]int, 0, 3) // for ambiguous components
		for i, name := range names {
			start, end := matchSpans[i*2], matchSpans[(i*2)+1]
			var sub string
			if start >= 0 && end >= 0 {
				sub = strings.ToLower(s[start:end])
			}

			switch name {
			case "year":
				year, e := strconv.Atoi(sub)
				if e == nil {
					year = ExtendYear(year)
					fd.SetYear(year)
				} else {
					fail = true
					break
				}
			case "month":
				month, e := strconv.Atoi(sub)
				if e == nil {
					// it was a number
					if month < 1 || month > 12 {
						fail = true
						break // month out of range
					}
					fd.SetMonth(month)
				} else {
					// try month name
					month, ok := monthLookup[sub]
					if !ok {
						fail = true
						break // nope.
					}
					fd.SetMonth(month)
				}
			case "cruftmonth":
				// special case to handle "Jan/Feb 2010"...
				// we'll make sure the first month is valid, then ignore it
				_, ok := monthLookup[sub]
				if !ok {
					fail = true
					break
				}
			case "day":
				day, e := strconv.Atoi(sub)
				if e != nil {
					fail = true
					break
				}
				if day < 1 || day > 31 {
					fail = true
					break
				}
				fd.SetDay(day)
			case "x1", "x2", "x3":
				// could be day, month or year...
				x, e := strconv.Atoi(sub)
				if e != nil {
					fail = true
					break
				}
				unknowns = append(unknowns, x)
			}
		}

		if fail {
			// regexp matched, but values sucked.
			continue
		}

		// got enough?
		if (fd.HasYear() && fd.HasMonth()) || (fd.HasMonth() && fd.HasDay()) {
			if fd.sane() {
				span.Begin, span.End = matchSpans[0], matchSpans[1]
				return fd, span, nil
			}
		} else {
			// got some ambiguous components to try?
			if len(unknowns) == 2 && fd.HasYear() {
				unknowns = append(unknowns, fd.Year())
			}
			if len(unknowns) == 3 {
				var err error
				fd, err = ctx.DateResolver(unknowns[0], unknowns[1], unknowns[2])
				if err != nil {
					return Date{}, Span{}, err
				}

				if fd.HasYear() && fd.HasMonth() && fd.HasDay() && fd.sane() {
					// resolved.
					span.Begin, span.End = matchSpans[0], matchSpans[1]
					return fd, span, nil
				}
			}
		}
	}

	// nothing. Just return an empty date and span
	return Date{}, Span{}, nil
}
