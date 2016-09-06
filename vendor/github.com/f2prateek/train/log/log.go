package log

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"sync"

	"github.com/f2prateek/train"
)

type Level uint8

const (
	// Logs nothing.
	None Level = iota
	// Logs request and response lines and their respective headers.
	Basic
	// Logs request and response lines and their respective headers and bodies (if present).
	Body
)

// New returns a logging interceptor with the given level that writes to the given writer.
func New(out io.Writer, level Level) train.Interceptor {
	return &loggingInterceptor{
		out:   out,
		level: level,
	}
}

type loggingInterceptor struct {
	out   io.Writer
	level Level
	sync.Mutex
}

func (interceptor *loggingInterceptor) Intercept(chain train.Chain) (*http.Response, error) {
	req := chain.Request()
	if interceptor.level == None {
		return chain.Proceed(req)
	}

	// Use a temp buffer so that a request/response pair always appears in order.
	var buf bytes.Buffer
	logBody := interceptor.level == Body

	requestDump, err := httputil.DumpRequestOut(req, logBody)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(&buf, "%s", requestDump)

	resp, err := chain.Proceed(req)

	responseDump, err := httputil.DumpResponse(resp, logBody)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(&buf, "%s", responseDump)

	interceptor.Lock()
	io.Copy(interceptor.out, &buf)
	interceptor.Unlock()

	return resp, err
}
