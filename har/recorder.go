package har

import (
	"fmt"
	"net/http"
	"time"

	martianhar "github.com/google/martian/har"
)

type Recorder struct {
	http.RoundTripper
	hl *martianhar.Logger
}

// NewRecorder returns a new Recorder object that fulfills the http.RoundTripper interface
func NewRecorder() *Recorder {
	hl := martianhar.NewLogger()
	hl.SetOption(martianhar.PostDataLogging(true))
	hl.SetOption(martianhar.BodyLogging(true))

	return &Recorder{
		RoundTripper: http.DefaultTransport,
		hl:           hl,
	}
}

func (r *Recorder) RoundTrip(req *http.Request) (*http.Response, error) {
	id := fmt.Sprintf("%d", time.Now().UnixNano())

	if err := r.hl.RecordRequest(id, req); err != nil {
		return nil, err
	}

	resp, err := r.RoundTripper.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	if err = r.hl.RecordResponse(id, resp); err != nil {
		return resp, err
	}

	return resp, err
}

// Export returns the in-memory log.
func (r *Recorder) Export() *martianhar.HAR {
	return r.hl.Export()
}
