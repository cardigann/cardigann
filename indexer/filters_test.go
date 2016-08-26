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
			t.Fatalf("Row %#d had an unexpected error: %s", idx+1, err.Error())
		}
		if result != example.expected {
			t.Fatalf("Row %#d was expecting %s, got %s", idx+1, example.expected, result)
		}
	}
}

func TestDateParseFilter(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	for idx, example := range []struct{ strTime, format, expected string }{
		{now.Format("Mon Jan 2 15:04:05 MST 2006"), "Mon Jan 2 15:04:05 MST 2006", now.Format(filterTimeFormat)},
	} {
		result, err := filterDateParse(example.format, example.strTime)
		if err != nil {
			t.Fatalf("Row %#d had an unexpected error: %s", idx+1, err.Error())
		}
		if result != example.expected {
			t.Fatalf("Row %#d was expecting %s, got %s", idx+1, example.expected, result)
		}
	}
}

func TestTimeAgoFilter(t *testing.T) {
	now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	for idx, example := range []struct {
		timeAgo  string
		expected time.Time
	}{
		{"0.5 months ago", now.Add(time.Hour * -372)},
		{"1 week, 2.5 days ago", now.Add((time.Hour * ((7 * 24) + 60)) * -1)},
		{"1 day ago", now.AddDate(0, 0, -1)},
		{"10.5 years", now.AddDate(-10, 0, -182).Add(time.Hour * -12)},
	} {
		result, err := filterTimeAgo(example.timeAgo, now)
		if err != nil {
			t.Fatalf("Row %#d had an unexpected error: %s", idx+1, err.Error())
		}
		if result != example.expected.Format(filterTimeFormat) {
			t.Fatalf("Row %#d was expecting %s, got %s", idx+1, example.expected.Format(filterTimeFormat), result)
		}
	}
}
