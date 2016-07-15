package bithdtv

import (
	"errors"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-humanize"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"

	"github.com/cardigann/cardigann/indexer"
	"github.com/cardigann/cardigann/torznab"
)

const (
	Key     = "bithdtv"
	baseURL = "https://www.bit-hdtv.com/"
)

var categoryMap = torznab.CategoryMapping{
	1:  torznab.Categories.TV_Anime,
	4:  torznab.Categories.TV_Documentary,
	5:  torznab.Categories.TV_Sport,
	10: torznab.Categories.TV,
	12: torznab.Categories.TV, // season packs
}

type torznabIndexer struct {
	config  indexer.Config
	browser browser.Browsable
}

func (i *torznabIndexer) login() error {
	err := i.browser.Open(baseURL + "login.php")
	if err != nil {
		return err
	}

	fm, err := i.browser.Form("form")
	if err != nil {
		return err
	}

	username, ok := i.config.Get(Key, "username")
	if !ok {
		return errors.New("No username in config for " + Key)
	}

	password, ok := i.config.Get(Key, "password")
	if !ok {
		return errors.New("No password in config for " + Key)
	}

	fm.Input("username", username)
	fm.Input("password", password)
	if fm.Submit() != nil {
		return err
	}

	// TODO: capture the following error message
	if strings.Contains(i.browser.Body(), "Login failed!") {
		return errors.New("Login failed, incorrect username or password")
	}

	log.Printf("Successfully logged in")
	return nil
}

func (i *torznabIndexer) Search(query torznab.Query) ([]torznab.ResultItem, error) {
	if err := i.login(); err != nil {
		return nil, err
	}

	localCat, hasLocalCat := query["cat"].(string)
	if !hasLocalCat {
		localCat = "0"
	}

	err := i.browser.OpenForm(baseURL+"torrents.php", url.Values{
		"search": []string{query["q"].(string)},
		"cat:":   []string{localCat},
	})
	if err != nil {
		return nil, err
	}

	items := []torznab.ResultItem{}

	i.browser.Find("table[width='750'] > tbody tr").Has("td.detail").Each(func(idx int, s *goquery.Selection) {
		link := s.Find("td").Eq(2).Find("a").First()
		title, ok := link.Attr("title")
		if !ok {
			log.Printf("Row %d has no title link", idx)
			return
		}

		downloadHref, ok := s.Find("td > p > a").First().Attr("href")
		if !ok {
			log.Printf("Row %d has no download link", idx)
			return
		}

		detailsHref, ok := s.Find("td").Eq(2).Find("a").First().Attr("href")
		if !ok {
			log.Printf("Row %d has no details link", idx)
			return
		}

		seeders, err := strconv.Atoi(s.Find("td").Eq(7).Text())
		if err != nil {
			log.Printf("Row %d has invalid seeders")
			return
		}

		leechers, err := strconv.Atoi(s.Find("td").Eq(8).Find("span").Text())
		if err != nil {
			log.Printf("Row %d has invalid leechers")
			return
		}

		catHref, ok := s.Find("td").Eq(1).Find("a").First().Attr("href")
		if err != nil {
			log.Printf("Row %d has invalid category href")
			return
		}

		catUrl, err := url.Parse(catHref)
		if err != nil {
			log.Println("Invalid category url", err)
			return
		}

		catId, err := strconv.Atoi(catUrl.Query().Get("cat"))
		if err != nil {
			log.Println("Invalid category id", err)
			return
		}

		mappedCat, ok := categoryMap[catId]
		if !ok {
			log.Printf("Failed to find a mapping for category %d", catId)
			return
		}

		bytes, err := humanize.ParseBytes(s.Find("td").Eq(6).Text())
		if err != nil {
			log.Println("Failed to parse file size", err)
			return
		}

		// timezone? assuming local
		pubTime, err := time.Parse("2006-01-02 15:04:05", s.Find("td").Eq(5).Text())
		if err != nil {
			log.Println("Failed to parse publish time", err)
			return
		}

		items = append(items, torznab.ResultItem{
			Title:           title,
			Description:     title,
			Comments:        detailsHref,
			GUID:            detailsHref,
			PublishDate:     pubTime,
			Link:            baseURL + downloadHref,
			Size:            bytes,
			Seeders:         seeders,
			Peers:           seeders + leechers,
			Category:        mappedCat.ID,
			MinimumRatio:    1,
			MinimumSeedTime: time.Hour * 48,
		})
	})

	return items, nil
}

func (i *torznabIndexer) Capabilities() torznab.Capabilities {
	return torznab.Capabilities{
		SearchModes: []torznab.SearchMode{
			{"search", true, []string{"q"}},
			{"tv-search", true, []string{"q", "season", "ep"}},
		},
	}
}

func init() {
	indexer.TorznabIndexers[Key] = func(c indexer.Config) (torznab.Indexer, error) {
		bow := surf.NewBrowser()
		bow.SetUserAgent(agent.Chrome())
		return &torznabIndexer{config: c, browser: bow}, nil
	}
}
