package bithdtv

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
	1:  torznab.CategoryTV_Anime,
	2:  torznab.CategoryMovies_BluRay,
	4:  torznab.CategoryTV_Documentary,
	5:  torznab.CategoryTV_Sport,
	6:  torznab.CategoryAudio,
	7:  torznab.CategoryMovies,
	8:  torznab.CategoryAudio_Video,
	10: torznab.CategoryTV,
	12: torznab.CategoryTV, // season packs
}

type torznabIndexer struct {
	config  indexer.Config
	browser browser.Browsable
}

func (i *torznabIndexer) Capabilities() torznab.Capabilities {
	return torznab.Capabilities{
		Categories: categoryMap,
		SearchModes: []torznab.SearchMode{
			{"search", true, []string{"q"}},
			{"tv-search", true, []string{"q", "season", "ep"}},
		},
	}
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

	log.Println(i.browser.Body())

	if strings.Contains(i.browser.Body(), "Login failed!") {
		msg := i.browser.Find("table.detail .text")
		msg.Find("style, b").Remove()
		return errors.New(strings.TrimSpace(msg.Text()))
	}

	log.Printf("Successfully logged in")
	return nil
}

func (i *torznabIndexer) Search(query torznab.Query) (*torznab.ResultFeed, error) {
	if err := i.login(); err != nil {
		return nil, err
	}

	catFilter, hasCatFilter := query["cat"].(string)

	keywords := query.Keywords()
	log.Printf("Searching for %q", keywords)

	err := i.browser.OpenForm(baseURL+"torrents.php", url.Values{
		"search": []string{keywords},
		"cat:":   []string{"0"},
	})
	if err != nil {
		return nil, err
	}

	items := []torznab.ResultItem{}
	timer := time.Now()

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

		// Search doesn't support multiple cats, so this filters
		if hasCatFilter && strconv.Itoa(mappedCat.ID) != catFilter {
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
			Site:            Key,
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

	log.Printf("Found %d results in %s", len(items), time.Now().Sub(timer))

	return &torznab.ResultFeed{
		Title:       "BIT-HDTV",
		Description: "Home of High Definition TV",
		Link:        "https://www.bit-hdtv.com/",
		Language:    "en-us",
		Items:       items,
	}, nil
}

func (i *torznabIndexer) Download(urlStr string) (io.ReadCloser, http.Header, error) {
	if err := i.login(); err != nil {
		return nil, http.Header{}, err
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, http.Header{}, err
	}

	u.Host = "www.bit-hdtv.com"
	u.Scheme = "https"

	if err := i.browser.Open(u.String()); err != nil {
		return nil, http.Header{}, err
	}

	b := &bytes.Buffer{}

	if _, err := i.browser.Download(b); err != nil {
		return nil, http.Header{}, err
	}

	return ioutil.NopCloser(bytes.NewReader(b.Bytes())), i.browser.ResponseHeaders(), nil
}

func init() {
	indexer.Registered[Key] = indexer.Constructor(func(c indexer.Config) (torznab.Indexer, error) {
		bow := surf.NewBrowser()
		bow.SetUserAgent(agent.Chrome())
		bow.SetAttribute(browser.SendReferer, false)
		bow.SetAttribute(browser.MetaRefreshHandling, false)
		return &torznabIndexer{config: c, browser: bow}, nil
	})
}
