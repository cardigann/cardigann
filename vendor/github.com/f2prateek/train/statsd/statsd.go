package statsd

import (
	"net/http"
	"time"

	"github.com/f2prateek/train"
	"github.com/statsd/client-interface"
)

// New returns a statsd interceptor that logs the following metrics:
// * requests (counter) - request count
// * requests.method.<method> (counter) - requests by method
// * request.size (timer) - request content-length
// * response.ok (counter) - successful requests
// * response.errors.client (counter) - client errors
// * response.errors.server (counter) - server errors
// * response.duration (timer) - request duration
// * response.size (timer) - response size in bytes
func New(stats statsd.Client) train.Interceptor {
	return &statsdInterceptor{stats}
}

type statsdInterceptor struct {
	statsd.Client
}

func (stats *statsdInterceptor) Intercept(chain train.Chain) (*http.Response, error) {
	req := chain.Request()
	stats.Incr("requests")
	stats.Incr("requests.method." + req.Method)
	stats.Histogram("request.size", int(req.ContentLength))

	start := time.Now()
	resp, err := chain.Proceed(req)
	elapsed := time.Now().Sub(start)

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		stats.Incr("response.ok")
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		stats.Incr("response.errors.client")
	case resp.StatusCode >= 500:
		stats.Incr("response.errors.server")
	}
	stats.Duration("response.duration", elapsed)
	stats.Histogram("response.size", int(resp.ContentLength))

	return resp, err
}
