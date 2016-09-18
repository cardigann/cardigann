package fuzzytime

import (
	"errors"
	"strings"
)

// DefaultContext is a predefined context which bails out if timezones or
// dates are ambiguous. It makes no attempt to resolve them.
var DefaultContext = Context{
	DateResolver: func(a, b, c int) (Date, error) {
		return Date{}, errors.New("ambiguous date")
	},
	TZResolver: DefaultTZResolver(""),
}

// USContext is a prefefined Context which opts for US timezones and mm/dd/yy dates
var USContext = Context{
	DateResolver: MDYResolver,
	TZResolver:   DefaultTZResolver("US"),
}

// WesternContext is a predefined Context which opts for UK and US timezones
// and dd/mm/yy dates
var WesternContext = Context{
	DateResolver: DMYResolver,
	TZResolver:   DefaultTZResolver("GB,US"),
}

// Extract tries to parse a Date and Time from a string.
// If none found, the returned DateTime will be empty
// Equivalent to DefaultContext.Extract()
func Extract(s string) (DateTime, []Span, error) { return DefaultContext.Extract(s) }

// ExtractTime tries to parse a Time from a string.
// Equivalent to DefaultContext.ExtractTime()
// Returns the parsed time information and a span showing which portion of the
// text matched, or an error.
func ExtractTime(s string) (Time, Span, error) { return DefaultContext.ExtractTime(s) }

// ExtractDate tries to parse a Date from a string.
// Equivalent to DefaultContext.ExtractDate()
// Returns the parsed date information and a span showing which portion of the
// text matched, or an error.
func ExtractDate(s string) (Date, Span, error) { return DefaultContext.ExtractDate(s) }

// Context provides helper functions to resolve ambiguous dates and timezones.
// For example, "CST" can mean China Standard Time, Central Standard Time in
// or Central Standard Time in Australia.
// Or, the date "5/2/10". It could Feburary 5th, 2010 or May 2nd 2010. Or even
// Feb 10th 2005, depending on country. Even "05/02/2010" is ambiguous.
// If you know something about the types of times and dates you're likely to
// encounter, you can provide a Context struct to guide the parsing.
type Context struct {
	// DateResolver is called when ambigous dates are encountered eg (10/11/12)
	// It should return a date, if one can be decided. Returning an error
	// indicates the resolver can't decide.
	DateResolver func(a, b, c int) (Date, error)
	// TZResolver returns the offset in seconds from UTC of the named zone (eg "EST").
	// if the resolver can't decide which timezone it is, it will return an error.
	TZResolver func(name string) (int, error)
}

// Extract tries to parse a Date and Time from a string
// It also returns a sorted list of spans specifing which bits of the
// string were used.
// If none found (or if there is an error), the returned
// DateTime will be empty.
func (ctx *Context) Extract(s string) (DateTime, []Span, error) {
	// do time first to cope with cases where the time breaks up the date: "Thu Aug 25 10:46:55 GMT 2011"
	ft, span1, err := ctx.ExtractTime(s)
	if err != nil {
		return DateTime{}, nil, err
	}
	if !ft.Empty() {
		// snip the matched time out of the string
		// (hack for nasty case where an hour can look like a 2-digit year)
		s = s[:span1.Begin] + s[span1.End:]
	}

	fd, span2, err := ctx.ExtractDate(s)
	if err != nil {
		return DateTime{}, nil, err
	}

	if !fd.Empty() {
		// fix up the second span to allow for the snipping
		if span2.Begin >= span1.Begin {
			span2.Begin += span1.End - span1.Begin
		}
		if span2.End >= span1.Begin {
			span2.End += span1.End - span1.Begin
		}
	}

	// sort/merge spans
	spans := tidySpans([]Span{span1, span2})

	return DateTime{fd, ft}, spans, nil
}

// DefaultTZResolver returns a TZResolver function which uses a list of country codes in
// preferredLocales to resolve ambigous timezones.
// For example, if you were expecting Bangladeshi times, then:
//     DefaultTZResolver("BD")
// would treat "BST" as Bangladesh Standard Time rather than British Summer Time
func DefaultTZResolver(preferredLocales string) func(name string) (int, error) {
	codes := strings.Split(strings.ToUpper(preferredLocales), ",")

	return func(name string) (int, error) {
		matches := FindTimeZone(name)
		if len(matches) == 1 {
			return TZToOffset(matches[0].Offset)
		} else if len(matches) > 1 {
			// try preferred locales in order of preference
			for _, cc := range codes {
				for _, tz := range matches {
					if strings.Contains(tz.Locale, cc) {
						return TZToOffset(tz.Offset)
					}
				}
			}
			return 0, errors.New("ambiguous timezone")
		} else {
			return 0, errors.New("unknown timezone")
		}
	}
}

// DMYResolver is a helper function for Contexts which treats
// ambiguous dates as DD/MM/YY
func DMYResolver(a, b, c int) (Date, error) {
	c = ExtendYear(c)
	return *NewDate(c, b, a), nil
}

// MDYResolver is a helper function for Contexts which treats
// ambiguous dates as MM/DD/YY
func MDYResolver(a, b, c int) (Date, error) {
	c = ExtendYear(c)
	return *NewDate(c, a, b), nil
}
