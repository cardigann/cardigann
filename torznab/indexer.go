package torznab

import (
	"io"
	"net/http"
)

type Info struct {
	ID          string
	Title       string
	Description string
	Link        string
	Language    string
	Category    string
}

type Indexer interface {
	Info() Info
	Search(query Query) ([]ResultItem, error)
	Download(urlStr string) (io.ReadCloser, http.Header, error)
	Capabilities() Capabilities
}
