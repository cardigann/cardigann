package tvmaze

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// DefaultClient is the default TV Maze client
var DefaultClient = NewClient()
var baseURL = url.URL{
	Scheme: "http",
	Host:   "api.tvmaze.com",
}

// Client represents a TV Maze client
type Client struct{}

// NewClient returns a new TV Maze client
func NewClient() Client {
	return Client{}
}

func (c Client) get(url url.URL, ret interface{}) (err error) {
	log.WithField("url", url.String()).Debug("getting url")
	r, err := http.Get(url.String())
	if err != nil {
		return err
	}

	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(&ret)
}

func baseURLWithPath(path string) url.URL {
	ret := baseURL
	ret.Path = path
	return ret
}

func baseURLWithPathQuery(path, key, val string) url.URL {
	ret := baseURL
	ret.Path = path
	ret.RawQuery = fmt.Sprintf("%s=%s", key, url.QueryEscape(val))
	return ret
}

func baseURLWithPathQueries(path string, vals map[string]string) url.URL {
	ret := baseURL
	ret.Path = path
	var queryStrings []string
	for key, val := range vals {
		queryStrings = append(queryStrings, fmt.Sprintf("%s=%s", key, url.QueryEscape(val)))
	}
	ret.RawQuery = strings.Join(queryStrings, "&")
	return ret
}
