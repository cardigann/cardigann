package torznab

import (
	"encoding/xml"
	"net/http"
	"sort"
	"strings"
)

type Capabilities struct {
	SearchModes []SearchMode
	Categories  CategoryMapping
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

	cats := c.Categories.Categories()
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
