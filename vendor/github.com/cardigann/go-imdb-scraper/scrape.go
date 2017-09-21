package imdbscraper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/headzoo/surf"
)

type Movie struct {
	Title string
	Year  string
}

func FindByID(id string) (*Movie, error) {
	bow := surf.NewBrowser()
	err := bow.Open(fmt.Sprintf("http://www.imdb.com/title/%s", id))
	if err != nil {
		return nil, err
	}

	title := bow.Dom().Find(".title_wrapper h1")
	if title.Length() == 0 {
		return nil, errors.New("Expected to find `.title_wrapper h1`")
	}

	year := title.Find("span#titleYear")
	if year.Length() == 0 {
		return nil, errors.New("Expected to find `span#titleYear` in title")
	}
	year.Remove()

	m := Movie{
		Year:  strings.Trim(strings.TrimSpace(year.Text()), "()"),
		Title: strings.TrimSpace(title.Text()),
	}

	return &m, nil
}
