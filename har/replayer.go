package har

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	martianhar "github.com/google/martian/har"
)

type Replayer struct {
	entries []*martianhar.Entry
}

func NewReplayerFromFile(path string) (*Replayer, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var h martianhar.HAR

	if err = json.Unmarshal(b, &h); err != nil {
		return nil, err
	}

	return &Replayer{entries: h.Log.Entries}, nil
}

func (r *Replayer) matchEntry(req *http.Request) (*martianhar.Entry, error) {
	if len(r.entries) == 0 {
		return nil, errors.New("No matching entry found")
	}

	entry := r.entries[0]
	r.entries = r.entries[1:]

	return entry, nil

}

func (r *Replayer) RoundTrip(req *http.Request) (*http.Response, error) {
	entry, err := r.matchEntry(req)
	if err != nil {
		panic(err)
	}
	resp, err := createResponse(entry.Response)
	if err != nil {
		panic(err)
	}
	resp.Request = req
	return resp, nil
}

func readHTTPVersion(v string) (string, int, int) {
	if v == "HTTP/1.0" {
		return "HTTP/1.0", 1, 0
	}
	return "HTTP/1.1", 1, 1
}

func createResponse(hresp *martianhar.Response) (*http.Response, error) {
	h := http.Header{}

	for _, hrow := range hresp.Headers {
		h.Add(hrow.Name, hrow.Value)
	}

	v, major, minor := readHTTPVersion(hresp.HTTPVersion)

	return &http.Response{
		Status:        hresp.StatusText,
		StatusCode:    hresp.Status,
		Header:        h,
		ContentLength: hresp.Content.Size,
		Body:          ioutil.NopCloser(bytes.NewReader(hresp.Content.Text)),
		Proto:         v,
		ProtoMajor:    major,
		ProtoMinor:    minor,
	}, nil
}
