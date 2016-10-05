package har_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cardigann/cardigann/har"
)

func TestHarRecorder(t *testing.T) {
	var responseStr = "Hello, client\n"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, responseStr)
	}))
	defer ts.Close()

	recorder := har.NewRecorder()
	c := http.Client{
		Transport: recorder,
	}

	res, err := c.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if string(greeting) != responseStr {
		t.Fatalf(`Expected %q, got %q`, responseStr, greeting)
	}

	h := recorder.Export()

	if got := len(h.Log.Entries); got != 1 {
		t.Fatalf("Expected 1 har entries, got %d", got)
	}
}
