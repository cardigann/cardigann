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
	err := bow.Open(fmt.Sprintf("http://www.imdb.com/title/%s/fullcredits", id))
	if err != nil {
		return nil, err
	}

	m := Movie{}
	m.Title = strings.TrimSpace(bow.Dom().Find("#main h3 > a").Text())
	m.Year = strings.Trim(strings.TrimSpace(bow.Dom().Find("#main h3 > span").Text()), "()")

	return &m, nil
}
