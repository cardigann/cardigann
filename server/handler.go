//go:generate esc -o static.go -prefix ../web/build -pkg server ../web/build
package server

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/cardigann/cardigann/indexer"
	"github.com/cardigann/cardigann/torznab"
	"github.com/gorilla/mux"

	// indexers
	_ "github.com/cardigann/cardigann/indexer/bithdtv"
)

const (
	buildDir         = "/web/build"
	localReactServer = "http://localhost:3000"
)

var (
	apiRoutePrefixes = []string{
		"/torznab/",
		"/download/",
		"/xhr/",
	}
)

type Params struct {
	BaseURL   string
	DevMode   bool
	SharedKey []byte
}

type handler struct {
	http.Handler
	Indexers    indexer.ConstructorMap
	Params      Params
	FileHandler http.Handler
}

func NewHandler(cm indexer.ConstructorMap, p Params) http.Handler {
	if p.DevMode {
		cm["example"] = indexer.Constructor(func(c indexer.Config) (torznab.Indexer, error) {
			return &indexer.ExampleIndexer{}, nil
		})
	}

	h := &handler{
		Indexers:    cm,
		Params:      p,
		FileHandler: http.FileServer(FS(false)),
	}

	if p.DevMode {
		u, err := url.Parse(localReactServer)
		if err != nil {
			panic(err)
		}

		log.Printf("Proxying static requests to %s", localReactServer)
		h.FileHandler = httputil.NewSingleHostReverseProxy(u)
	}

	router := mux.NewRouter()

	// torznab routes
	router.HandleFunc("/torznab/{indexer}/api", h.torznabHandler).Methods("GET")
	router.HandleFunc("/download/{token}/{filename}", h.downloadHandler).Methods("GET")

	// xhr routes for the webapp
	xhr := xhrHandler{}
	router.HandleFunc("/xhr/indexers", h.getIndexersHandler).Methods("GET")

	h.Handler = router
	return h
}

func (h *handler) baseURL(r *http.Request, path string) (*url.URL, error) {
	if h.Params.BaseURL != "" {
		return url.Parse(h.Params.BaseURL)
	}
	proto := "http"
	if r.TLS != nil {
		proto = "https"
	}
	return url.Parse(fmt.Sprintf("%s://%s%s", proto, r.Host, path))
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	for _, prefix := range apiRoutePrefixes {
		if strings.HasPrefix(r.URL.Path, prefix) {
			h.APIHandler.ServeHTTP(w, r)
			return
		}
	}

	h.FileHandler.ServeHTTP(w, r)
}

func (h *handler) torznabHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	indexerID := params["indexer"]

	indexer, err := indexer.Registered.New(indexerID)
	if err != nil {
		log.Fatal(err)
	}

	t := r.URL.Query().Get("t")

	switch t {
	case "caps":
		indexer.Capabilities().ServeHTTP(w, r)

	case "search", "tvsearch", "tv-search":
		feed, err := h.search(r, indexer, indexerID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		x, err := xml.MarshalIndent(feed, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		w.Write(x)

	default:
		http.Error(w, fmt.Sprintf("Type %q not implemented", t), http.StatusNotFound)
	}
}

func (h *handler) downloadHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	token := params["token"]
	filename := params["filename"]

	t, err := decodeToken(token, h.Params.SharedKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	indexer, err := indexer.Registered.New(t.Site)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	rc, headers, err := indexer.Download(t.Link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Length", headers.Get("Content-Length"))
	w.Header().Set("Content-Type", "application/x-download")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Transfer-Encoding", "binary")

	defer rc.Close()
	go io.Copy(w, rc)
}

func (h *handler) search(r *http.Request, indexer torznab.Indexer, siteKey string) (*torznab.ResultFeed, error) {
	baseURL, err := h.baseURL(r, "/download")
	if err != nil {
		return nil, err
	}

	query, err := torznab.ParseQuery(r.URL.Query())
	if err != nil {
		return nil, err
	}

	log.Printf("Query: %#v", query)
	feed, err := indexer.Search(query)
	if err != nil {
		return nil, err
	}

	// rewrite links to use the server
	for idx, item := range feed.Items {
		t := &token{
			Site: item.Site,
			Link: item.Link,
		}
		te, err := t.Encode(h.Params.SharedKey)
		if err != nil {
			return nil, err
		}
		baseURL.Path += fmt.Sprintf("/%s/%s.torrent", te, item.Title)
		feed.Items[idx].Link = baseURL.String()
	}

	return feed, err
}
