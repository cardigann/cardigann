package train_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/f2prateek/train"
	"github.com/f2prateek/train/log"
	"github.com/f2prateek/train/mocks"
	"github.com/gohttp/response"
	"github.com/stretchr/testify/mock"
)

// Interceptor that simply calls `chain.Proceed(chain.Request())`.
var fallThrough = train.InterceptorFunc(func(chain train.Chain) (*http.Response, error) {
	return chain.Proceed(chain.Request())
})

func TestFallThroughInterceptor(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		response.OK(w, "Hello World!")
	}))
	defer ts.Close()

	client := &http.Client{
		Transport: train.Transport(fallThrough),
	}

	resp, err := client.Get(ts.URL)
	assert.Equal(t, nil, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Hello World!\n", string(body))
}

func TestInterceptorCanShortCircuit(t *testing.T) {
	shortCircuitInterceptor := mocks.New()
	{
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			response.OK(w, "Hello World!")
		}))
		defer ts.Close()
		resp, err := http.Get(ts.URL)
		shortCircuitInterceptor.On("Intercept", mock.AnythingOfType("*train.interceptorChain")).Return(resp, err)
	}
	m := mocks.New()

	client := &http.Client{
		Transport: train.Transport(fallThrough, shortCircuitInterceptor, m),
	}

	resp, err := client.Get("https://golang.org/")

	// Assert that the application sees our "shorted" response.
	assert.Equal(t, nil, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Hello World!\n", string(body))

	// Assert our mocks.
	shortCircuitInterceptor.AssertExpectations(t)
	m.AssertExpectations(t)
}

func ExampleTransport() {
	client := &http.Client{
		// Try changing the log level!
		Transport: train.Transport(log.New(os.Stdout, log.None)),
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		response.OK(w, "Hello World!")
	}))
	defer ts.Close()

	resp, _ := client.Get(ts.URL)
	bytes, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(bytes))
	// Output: Hello World!
}
