package torznab

import (
	"io"
	"net/url"
)

type Indexer interface {
	Search(query Query) (*ResultFeed, error)
	Download(u *url.URL) (io.ReadCloser, error)
	Capabilities() Capabilities
}
