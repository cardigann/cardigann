package indexer

import (
	"testing"
	"time"
)

func TestQueryStringFilter(t *testing.T) {
	for idx, example := range []struct{ u, param, expected string }{
		{"http://example.org/test?llamas=1", "llamas", "1"},
		{"http://example.org/test", "llamas", ""},
		{"/test?llamas=true", "llamas", "true"},
		{"test?llamas=true", "llamas", "true"},
	} {
		result, err := filterQueryString(example.param, example.u)
		if err != nil {
			t.Fatalf("Row #%d had an unexpected error: %s", idx+1, err.Error())
		}
		if result != example.expected {
			t.Fatalf("Row #%d was expecting %s, got %s", idx+1, example.expected, result)
		}
	}
}

func TestDateParseFilter(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	for idx, example := range []struct{ strTime, format, expected string }{
		{now.Format("Mon Jan 2 15:04:05 MST 2006"), "Mon Jan 2 15:04:05 MST 2006", now.Format(filterTimeFormat)},
	} {
		result, err := filterDateParse([]string{example.format}, example.strTime)
		if err != nil {
			t.Fatalf("Row #%d had an unexpected error: %s", idx+1, err.Error())
		}
		if result != example.expected {
			t.Fatalf("Row #%d was expecting %s, got %s", idx+1, example.expected, result)
		}
	}
}

func TestFuzzyTimeFilter(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	for idx, example := range []struct {
		pattern  string
		expected time.Time
	}{
		{"2009-11-10T23:00:00Z", now},
		{"Today 10:00am", now.Add(time.Hour * -13)},
		{"Today, 10:00am", now.Add(time.Hour * -13)},
		{"Today", now.Add(time.Hour * -23)},
		{"Yesterday", now.AddDate(0, 0, -1).Add(time.Hour * -23)},
		{"Yesterday 10:00am", now.AddDate(0, 0, -1).Add(time.Hour * -13)},
		{"3 mins ago", now.Add(time.Minute * -3)},
		{"06-01 19:39", time.Date(2009, time.June, 01, 19, 39, 0, 0, time.UTC)},
		{"06-01-2009 19:39", time.Date(2009, time.June, 01, 19, 39, 0, 0, time.UTC)},
		{"06-01-09 19:39", time.Date(2009, time.June, 01, 19, 39, 0, 0, time.UTC)},
	} {
		result, err := filterFuzzyTime(example.pattern, now)
		if err != nil {
			t.Fatalf("Row #%d had an unexpected error: %s", idx+1, err.Error())
		}
		if result != example.expected.Format(filterTimeFormat) {
			t.Fatalf("Row #%d was expecting %s, got %s", idx+1, example.expected.Format(filterTimeFormat), result)
		}
	}
}

func TestParseTimeAgo(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	for idx, example := range []struct {
		timeAgo  string
		expected time.Time
	}{
		{"0.5 months ago", now.Add(time.Hour * -372)},
		{"1 week, 2.5 days ago", now.Add((time.Hour * ((7 * 24) + 60)) * -1)},
		{"1 day ago", now.AddDate(0, 0, -1)},
		{"10.5 years", now.AddDate(-10, 0, -182).Add(time.Hour * -12)},
		{"3 mins ago", now.Add(time.Minute * -3)},
		{"57m 7s ago", now.Add(time.Minute * -57).Add(time.Second * -7)},
	} {
		result, err := parseTimeAgo(example.timeAgo, now)
		if err != nil {
			t.Fatalf("Row #%d had an unexpected error: %s", idx+1, err.Error())
		}
		if !result.Equal(example.expected) {
			t.Fatalf("Row #%d was expecting %s, got %s",
				idx+1, example.expected.Format(filterTimeFormat), result.Format(filterTimeFormat))
		}
	}
}
