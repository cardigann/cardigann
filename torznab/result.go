package torznab

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
)

const rfc822 = "Mon, 02 Jan 2006 15:04:05 -0700"

type ResultItem struct {
	Site        string
	Title       string
	Description string
	GUID        string
	Comments    string
	Link        string
	Category    int
	Size        uint64
	Files       int
	Grabs       int
	PublishDate time.Time

	Seeders              int
	Peers                int
	MinimumRatio         float64
	MinimumSeedTime      time.Duration
	DownloadVolumeFactor float64
	UploadVolumeFactor   float64
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
		Files       int         `xml:"files,omitempty"`
		Grabs       int         `xml:"grabs,omitempty"`
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
		Files:       ri.Files,
		Grabs:       ri.Grabs,
		PublishDate: ri.PublishDate.Format(rfc822),
		Enclosure:   enclosure,
		Attrs: []torznabAttrView{
			{Name: "site", Value: ri.Site},
			{Name: "seeders", Value: strconv.Itoa(ri.Seeders)},
			{Name: "peers", Value: strconv.Itoa(ri.Peers)},
			{Name: "minimumratio", Value: fmt.Sprintf("%.2f", ri.MinimumRatio)},
			{Name: "minimumseedtime", Value: fmt.Sprintf("%.f", ri.MinimumSeedTime.Seconds())},
			{Name: "size", Value: fmt.Sprintf("%d", ri.Size)},
			{Name: "downloadvolumefactor", Value: fmt.Sprintf("%.2f", ri.DownloadVolumeFactor)},
			{Name: "uploadvolumefactor", Value: fmt.Sprintf("%.2f", ri.UploadVolumeFactor)},
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
	Info  Info
	Items []ResultItem
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
		Title:       rf.Info.Title,
		Description: rf.Info.Description,
		Link:        rf.Info.Link,
		Language:    rf.Info.Language,
		Category:    rf.Info.Category,
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
