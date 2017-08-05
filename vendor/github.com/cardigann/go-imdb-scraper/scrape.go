package imdbscraper

import (
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

	m := Movie{}
	m.Year = strings.Trim(strings.TrimSpace(bow.Dom().Find(".title_wrapper #titleYear a").Text()), "()")

	bow.Dom().Find(".title_wrapper h1 *").Remove()
	m.Title = strings.TrimSpace(bow.Dom().Find(".title_wrapper h1").Text())

	return &m, nil
}
