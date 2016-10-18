package torznab

import (
	"encoding/xml"
	"net/http"
	"sort"
	"strings"
)

type Capabilities struct {
	SearchModes []SearchMode
	Categories  Categories
}

func (c Capabilities) HasSearchMode(key string) (bool, []string) {
	for _, m := range c.SearchModes {
		if m.Key == key && m.Available {
			return true, m.SupportedParams
		}
	}
	return false, nil
}

func (c Capabilities) HasTVShows() bool {
	for _, cat := range c.Categories {
		if cat.ID >= 5000 && cat.ID < 6000 {
			return true
		}
	}
	return false
}

func (c Capabilities) HasMovies() bool {
	for _, cat := range c.Categories {
		if cat.ID >= 2000 && cat.ID < 3000 {
			return true
		}
	}
	return false
}

type SearchMode struct {
	Key             string
	Available       bool
	SupportedParams []string
}

func (c Capabilities) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var cx struct {
		XMLName   struct{} `xml:"caps"`
		Searching struct {
			Values []interface{}
		} `xml:"searching"`
		Categories struct {
			Values []interface{}
		} `xml:"categories"`
	}

	for _, mode := range c.SearchModes {
		available := "no"
		if mode.Available {
			available = "yes"
		}
		cx.Searching.Values = append(cx.Searching.Values, struct {
			XMLName         xml.Name
			Available       string `xml:"available,attr"`
			SupportedParams string `xml:"supportedParams,attr"`
		}{
			xml.Name{"", mode.Key},
			available,
			strings.Join(mode.SupportedParams, ","),
		})
	}

	cats := c.Categories
	sort.Sort(cats)

	for _, cat := range cats {
		cx.Categories.Values = append(cx.Categories.Values, struct {
			XMLName struct{} `xml:"category"`
			ID      int      `xml:"id,attr"`
			Name    string   `xml:"name,attr"`
		}{
			ID:   cat.ID,
			Name: cat.Name,
		})
	}

	e.Encode(cx)
	return nil
}

func (c Capabilities) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	x, err := xml.MarshalIndent(c, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}
