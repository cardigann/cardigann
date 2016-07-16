package torznab

import (
	"encoding/xml"
	"net/http"
	"strings"
)

type Capabilities struct {
	SearchModes []SearchMode
}

type SearchMode struct {
	Key             string
	Available       bool
	SupportedParams []string
}

func (c Capabilities) HasSearchMode(t string) bool {
	for _, m := range c.SearchModes {
		if m.Key == t {
			return true
		}
	}
	return false
}

func (c Capabilities) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var cx struct {
		XMLName struct{}      `xml:"caps"`
		Modes   []interface{} `xml:"searching"`
	}

	for _, mode := range c.SearchModes {
		available := "no"
		if mode.Available {
			available = "yes"
		}
		cx.Modes = append(cx.Modes, struct {
			XMLName         xml.Name
			Available       string `xml:"available,attr"`
			SupportedParams string `xml:"supportedParams,attr"`
		}{
			xml.Name{"", mode.Key},
			available,
			strings.Join(mode.SupportedParams, ","),
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
