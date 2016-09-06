package curl_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/f2prateek/train"
	"github.com/f2prateek/train/curl"
	"github.com/gohttp/response"
)

func TestCurl(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.OK(w)
	}))
	defer ts.Close()

	cases := []struct {
		Method  string
		Body    io.Reader
		Headers map[string]string
		Output  string // ts.URL will automatically be appended.
	}{
		{"GET", nil, nil, "curl -X GET"},
		{"POST", strings.NewReader("foo"), nil, "curl -X POST -d foo"},
		{"PUT", nil, map[string]string{"foo": "bar"}, "curl -X PUT -H \"Foo: bar\""},
	}

	for _, c := range cases {
		buf := new(bytes.Buffer)
		client := &http.Client{
			Transport: train.Transport(curl.New(buf)),
		}

		req, err := http.NewRequest(c.Method, ts.URL, c.Body)
		for k, v := range c.Headers {
			req.Header.Add(k, v)
		}
		assert.Equal(t, nil, err)

		resp, err := client.Do(req)
		assert.Equal(t, nil, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, c.Output+" "+ts.URL, buf.String())
	}
}
