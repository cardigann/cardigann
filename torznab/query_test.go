package torznab

import (
	"net/url"
	"testing"
)

func TestParsingQueryKeywords(t *testing.T) {
	var rows = []struct {
		Vals     url.Values
		Expected string
	}{
		{url.Values{"q": []string{"llamas"}}, "llamas"},
		{url.Values{"q": []string{"llamas"}, "season": []string{"2"}, "ep": []string{"12"}}, "llamas S02E12"},
	}

	for idx, row := range rows {
		q, err := ParseQuery(row.Vals)
		if err != nil {
			t.Fatal(err)
		}

		k := q.Keywords()
		if k != row.Expected {
			t.Fatalf(`Row %d: Expected q=%q, got q=%q`, idx+1, row.Expected, k)
		}
	}
}
