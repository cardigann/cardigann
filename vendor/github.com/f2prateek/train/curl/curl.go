package curl

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/f2prateek/train"
)

// New returns an interceptor interceptor that logs request as curl shell
// commands to the given io.Writer.
func New(out io.Writer) train.Interceptor {
	return &curlInterceptor{
		out: out,
	}
}

type curlInterceptor struct {
	out io.Writer
	sync.Mutex
}

func (interceptor *curlInterceptor) Intercept(chain train.Chain) (*http.Response, error) {
	req := chain.Request()

	cmd := "curl"
	cmd = cmd + " -X " + req.Method
	for k, values := range req.Header {
		for _, v := range values {
			cmd = cmd + " -H " + "\"" + k + ": " + v + "\""
		}
	}

	if req.Body != nil {
		// Copy the original body into a buffer.
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(req.Body); err != nil {
			return nil, err
		}
		if err := req.Body.Close(); err != nil {
			return nil, err
		}

		// Replace the request body with copy of the buffer.
		req.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))

		// Log the body.
		cmd = cmd + " -d " + buf.String()
	}

	cmd = cmd + " " + req.URL.String()

	interceptor.Lock()
	io.Copy(interceptor.out, strings.NewReader(cmd))
	interceptor.Unlock()

	return chain.Proceed(req)
}
