package log_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	"github.com/f2prateek/train"
	"github.com/f2prateek/train/log"
	"github.com/gohttp/response"
)

func TestLevels(t *testing.T) {
	cases := []struct {
		level log.Level
		out   string
	}{
		{log.None, ""},
		{log.Basic, `GET / HTTP/1.1
Host: 127.0.0.1:54709
User-Agent: Go-http-client/1.1
Accept-Encoding: gzip

HTTP/1.1 200 OK
Content-Length: 13
Content-Type: text/plain; charset=utf-8
Date: Tue, 10 Nov 2009 23:00:00 UTC

`},
		{log.Body, `GET / HTTP/1.1
Host: 127.0.0.1:54709
User-Agent: Go-http-client/1.1
Accept-Encoding: gzip

HTTP/1.1 200 OK
Content-Length: 13
Content-Type: text/plain; charset=utf-8
Date: Tue, 10 Nov 2009 23:00:00 UTC

Hello World!
`},
	}

	for _, c := range cases {
		var buf bytes.Buffer
		client := &http.Client{
			Transport: train.Transport(log.New(&buf, c.level)),
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Date", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC).Format(time.RFC1123))
			response.OK(w, "Hello World!")
		}))
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			panic(err)
		}
		req.Host = "127.0.0.1:54709"
		client.Do(req)

		assert.Equal(t, c.out, strings.Replace(buf.String(), "\r", "", -1))
	}
}

func ExampleNew() {
	var buf bytes.Buffer
	client := &http.Client{
		Transport: train.Transport(log.New(&buf, log.Body)),
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Date", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC).Format(time.RFC1123))
		response.OK(w, "Hello World!")
	}))
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		panic(err)
	}
	req.Host = "127.0.0.1:54709"
	client.Do(req)

	fmt.Println(strings.Replace(buf.String(), "\r", "", -1))
	// Output:
	// GET / HTTP/1.1
	// Host: 127.0.0.1:54709
	// User-Agent: Go-http-client/1.1
	// Accept-Encoding: gzip
	//
	// HTTP/1.1 200 OK
	// Content-Length: 13
	// Content-Type: text/plain; charset=utf-8
	// Date: Tue, 10 Nov 2009 23:00:00 UTC
	//
	// Hello World!
}
