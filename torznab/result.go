package torznab

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
)

const rfc2822 = "Mon Jan 02 2006 15:04:05 -0700"

type ResultItem struct {
	Title       string
	Description string
	GUID        string
	Comments    string
	Link        string
	Category    int
	Size        uint64
	PublishDate time.Time

	Seeders         int
	Peers           int
	MinimumRatio    float64
	MinimumSeedTime time.Duration
}

func (ri ResultItem) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var enclosure = struct {
		URL    string `xml:"url,attr,omitempty"`
		Length uint64 `xml:"length,attr,omitempty"`
		Type   string `xml:"type,attr,omitempty"`
	}{
		URL:    ri.Link,
		Length: ri.Size,
		Type:   "application/x-bittorrent",
	}

	var itemView = struct {
		XMLName struct{} `xml:"item"`

		// standard rss elements
		Title       string      `xml:"title,omitempty"`
		Description string      `xml:"description,omitempty"`
		GUID        string      `xml:"guid,omitempty"`
		Comments    string      `xml:"comments,omitempty"`
		Link        string      `xml:"link,omitempty"`
		Category    string      `xml:"category,omitempty"`
		Size        uint64      `xml:"size,omitempty"`
		PublishDate string      `xml:"pubDate,omitempty"`
		Enclosure   interface{} `xml:"enclosure,omitempty"`

		// torznab elements
		Attrs []torznabAttrView
	}{
		Title:       ri.Title,
		Description: ri.Description,
		GUID:        ri.GUID,
		Comments:    ri.Comments,
		Link:        ri.Link,
		Category:    strconv.Itoa(ri.Category),
		Size:        ri.Size,
		PublishDate: ri.PublishDate.Format(rfc2822),
		Enclosure:   enclosure,
		Attrs: []torznabAttrView{
			{Name: "seeders", Value: strconv.Itoa(ri.Seeders)},
			{Name: "peers", Value: strconv.Itoa(ri.Peers)},
			{Name: "minimumratio", Value: fmt.Sprintf("%.f", ri.MinimumRatio)},
			{Name: "minimumseedtime", Value: fmt.Sprintf("%.f", ri.MinimumSeedTime.Seconds())},
		},
	}

	e.Encode(itemView)
	return nil
}

type torznabAttrView struct {
	XMLName struct{} `xml:"torznab:attr"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

type ResultFeed struct {
	Title       string
	Description string
	Link        string
	Language    string
	Category    string
	Items       []ResultItem
}

func (rf ResultFeed) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var channelView = struct {
		XMLName     struct{} `xml:"channel"`
		Title       string   `xml:"title,omitempty"`
		Description string   `xml:"description,omitempty"`
		Link        string   `xml:"link,omitempty"`
		Language    string   `xml:"language,omitempty"`
		Category    string   `xml:"category,omitempty"`
		Items       []ResultItem
	}{
		Title:       rf.Title,
		Description: rf.Description,
		Link:        rf.Link,
		Language:    rf.Language,
		Category:    rf.Category,
		Items:       rf.Items,
	}

	e.Encode(struct {
		XMLName          struct{}    `xml:"rss"`
		TorznabNamespace string      `xml:"xmlns:torznab,attr"`
		Version          string      `xml:"version,attr,omitempty"`
		Channel          interface{} `xml:"channel"`
	}{
		Version:          "2.0",
		Channel:          channelView,
		TorznabNamespace: "http://torznab.com/schemas/2015/feed",
	})
	return nil
}
