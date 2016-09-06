// Package train provides a http.RoundTripper with chainable middleware.
package train

import "net/http"

type Chain interface {
	// Request returns the http.Request for this chain.
	Request() *http.Request
	// Proceed the chain with a given request and returns the result.
	Proceed(*http.Request) (*http.Response, error)
}

// Observes, modifies, and potentially short-circuits requests going out and the corresponding
// requests coming back in. Typically interceptors will be used to add, remove, or transform headers
// on the request or response. Interceptors must return either a response or an error.
type Interceptor interface {
	// Intercept the chain and return a result.
	Intercept(Chain) (*http.Response, error)
}

// The InterceptorFunc type is an adapter to allow the use of ordinary functions as interceptors.
// If f is a function with the appropriate signature, InterceptorFunc(f) is a Interceptor that calls f.
type InterceptorFunc func(Chain) (*http.Response, error)

// Intercept calls f(c).
func (f InterceptorFunc) Intercept(c Chain) (*http.Response, error) {
	return f(c)
}

// Return a new http.RoundTripper with the given interceptors and http.DefaultTransport.
// Interceptors will be called in the order they are provided.
func Transport(interceptors ...Interceptor) http.RoundTripper {
	return TransportWith(http.DefaultTransport, interceptors...)
}

// Return a new http.RoundTripper with the given interceptors and a custom http.RoundTripper
// to perform the actual HTTP request. Interceptors will be called in the order they are
// provided.
func TransportWith(transport http.RoundTripper, interceptors ...Interceptor) http.RoundTripper {
	return &interceptorRoundTripper{
		interceptors: append([]Interceptor{}, interceptors...),
		transport:    transport,
	}
}

type interceptorRoundTripper struct {
	interceptors []Interceptor
	transport    http.RoundTripper
}

func (i *interceptorRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	chain := &interceptorChain{
		index:        0,
		request:      req,
		interceptors: i.interceptors,
		transport:    i.transport,
	}
	return chain.Proceed(req)
}

type interceptorChain struct {
	index        int
	request      *http.Request
	interceptors []Interceptor
	transport    http.RoundTripper
}

func (c *interceptorChain) Request() *http.Request {
	return c.request
}

func (c *interceptorChain) Proceed(req *http.Request) (*http.Response, error) {
	// If there's another interceptor in the chain, call that.
	if c.index < len(c.interceptors) {
		chain := &interceptorChain{
			index:        c.index + 1,
			request:      req,
			interceptors: c.interceptors,
			transport:    c.transport,
		}
		interceptor := c.interceptors[c.index]
		return interceptor.Intercept(chain)
	}

	// No more interceptors. Do HTTP.
	return c.transport.RoundTrip(req)
}
