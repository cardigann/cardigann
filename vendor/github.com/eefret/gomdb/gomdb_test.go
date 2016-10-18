package gomdb

import "testing"

func TestSearch(t *testing.T) {
	tests := []struct {
		query *QueryData
		title string
		year  string
	}{
		{&QueryData{Title: "Fight Club", Year: "1999", SearchType: MovieSearch},
			"Fight Club",
			"1999",
		},
		{&QueryData{Title: "Her"},
			"Her",
			"2013",
		},
		{&QueryData{Title: "Macbeth", Year: "2015"},
			"Macbeth",
			"2015",
		},
	}

	for i, item := range tests {
		resp, err := Search(item.query)
		if err != nil {
			t.Errorf("Test[%d]: %s", i, err)
			continue
		}
		if resp.Search[0].Title != item.title {
			t.Errorf("Test[%d]: Expected- %s, Got- %s", i, item.title, resp.Search[0].Title)
			continue
		}
		if resp.Search[0].Year != item.year {
			t.Errorf("Test[%d]: Expected- %s, Got- %s", i, item.year, resp.Search[0].Year)
			continue
		}
	}
}

func TestFailSearch(t *testing.T) {
	tests := []struct {
		query *QueryData
	}{
		{&QueryData{Title: "Game of Thrones", Year: "2001"}},
		{&QueryData{Title: "Dexter", SearchType: EpisodeSearch}},
	}

	for i, item := range tests {
		_, err := Search(item.query)
		if err == nil {
			t.Errorf("Test[%d]: Got nil error", i)
			continue
		}
		// Checking for strings is bad. But the API might change.
		if err.Error() != "Movie not found!" {
			t.Errorf("Test[%d]: Unexpected value- %s, Got- %s", i, err)
			continue
		}
	}
}

func TestInvalidCategory(t *testing.T) {
	tests := []struct {
		query *QueryData
	}{
		{&QueryData{Title: "Game of Thrones", Year: "2001", SearchType: "bad"}},
		{&QueryData{Title: "Dexter", SearchType: "bad"}},
	}

	for i, item := range tests {
		_, err := Search(item.query)
		if err == nil {
			t.Errorf("Test[%d]: Got nil error", i)
			continue
		}
		// Checking for strings is bad. But the error type is formatted
		if err.Error() != "Invalid search category- bad" {
			t.Errorf("Test[%d]: Unexpected value- %s, Got- %s", i, err)
			continue
		}
	}
}

func TestMovieByTitle(t *testing.T) {
	tests := []struct {
		query *QueryData
		title string
		year  string
	}{
		{&QueryData{Title: "Fight Club", Year: "1999", SearchType: MovieSearch},
			"Fight Club",
			"1999",
		},
		{&QueryData{Title: "Her"},
			"Her",
			"2013",
		},
		{&QueryData{Title: "Macbeth", Year: "2015"},
			"Macbeth",
			"2015",
		},
	}

	for i, item := range tests {
		resp, err := MovieByTitle(item.query)
		if err != nil {
			t.Errorf("Test[%d]: %s", i, err)
			continue
		}
		if resp.Title != item.title {
			t.Errorf("Test[%d]: Expected- %s, Got- %s", i, item.title, resp.Title)
			continue
		}
		if resp.Year != item.year {
			t.Errorf("Test[%d]: Expected- %s, Got- %s", i, item.year, resp.Year)
			continue
		}
	}
}

func TestMovieByImdbID(t *testing.T) {
	tests := []struct {
		id    string
		title string
		year  string
	}{
		{
			"tt0137523",
			"Fight Club",
			"1999",
		},
		{
			"tt1798709",
			"Her",
			"2013",
		},
		{
			"tt2884018",
			"Macbeth",
			"2015",
		},
	}

	for i, item := range tests {
		resp, err := MovieByImdbID(item.id)
		if err != nil {
			t.Errorf("Test[%d]: %s", i, err)
			continue
		}
		if resp.Title != item.title {
			t.Errorf("Test[%d]: Expected- %s, Got- %s", i, item.title, resp.Title)
			continue
		}
		if resp.Year != item.year {
			t.Errorf("Test[%d]: Expected- %s, Got- %s", i, item.year, resp.Year)
			continue
		}
	}
}
