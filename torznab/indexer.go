package torznab

import (
	"io"
	"net/http"
)

type Indexer interface {
	Search(query Query) (*ResultFeed, error)
	Download(urlStr string) (io.ReadCloser, http.Header, error)
	Capabilities() Capabilities
}
