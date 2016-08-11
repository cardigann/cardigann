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
	BaseURL    string
	DevMode    bool
	APIKey     []byte
	Passphrase string
	Config     indexer.Config
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

		k, _ := h.sharedKey()
		log.Printf("API Key is %x", k)
	}

	router := mux.NewRouter()

	// torznab routes
	router.HandleFunc("/torznab/{indexer}", h.torznabHandler).Methods("GET")
	router.HandleFunc("/torznab/{indexer}/api", h.torznabHandler).Methods("GET")
	router.HandleFunc("/download/{token}/{filename}", h.downloadHandler).Methods("GET")

	// xhr routes for the webapp
	router.HandleFunc("/xhr/indexers/{indexer}/test", h.postIndexerTestHandler).Methods("POST")
	router.HandleFunc("/xhr/indexers/{indexer}/config", h.getIndexersConfigHandler).Methods("GET")
	router.HandleFunc("/xhr/indexers/{indexer}/config", h.patchIndexersConfigHandler).Methods("PATCH")
	router.HandleFunc("/xhr/indexers", h.getIndexersHandler).Methods("GET")
	router.HandleFunc("/xhr/indexers", h.patchIndexersHandler).Methods("PATCH")
	router.HandleFunc("/xhr/auth", h.postAuthHandler).Methods("POST")

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
	log.Printf("%s %s", r.Method, r.URL.RequestURI())

	for _, prefix := range apiRoutePrefixes {
		if strings.HasPrefix(r.URL.Path, prefix) {
			h.Handler.ServeHTTP(w, r)
			return
		}
	}

	h.FileHandler.ServeHTTP(w, r)
}

func (h *handler) torznabHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	indexerID := params["indexer"]

	apiKey := r.URL.Query().Get("apikey")
	if !h.checkAPIKey(apiKey) {
		torznab.Error(w, "Invalid apikey parameter", torznab.ErrInsufficientPrivs)
		return
	}

	indexer, err := indexer.Registered.New(indexerID, h.Params.Config)
	if err != nil {
		torznab.Error(w, err.Error(), torznab.ErrIncorrectParameter)
		return
	}

	t := r.URL.Query().Get("t")

	if t == "" {
		http.Redirect(w, r, r.URL.Path+"?t=caps", http.StatusTemporaryRedirect)
		return
	}

	switch t {
	case "caps":
		indexer.Capabilities().ServeHTTP(w, r)

	case "search", "tvsearch", "tv-search":
		feed, err := h.search(r, indexer, indexerID)
		if err != nil {
			torznab.Error(w, err.Error(), torznab.ErrUnknownError)
			return
		}
		x, err := xml.MarshalIndent(feed, "", "  ")
		if err != nil {
			torznab.Error(w, err.Error(), torznab.ErrUnknownError)
			return
		}
		w.Header().Set("Content-Type", "application/rss+xml")
		w.Write(x)

	default:
		torznab.Error(w, "Unknown type parameter", torznab.ErrIncorrectParameter)
	}
}

func (h *handler) downloadHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	token := params["token"]
	filename := params["filename"]

	k, err := h.sharedKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := decodeToken(token, k)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	indexer, err := indexer.Registered.New(t.Site, h.Params.Config)
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
	items, err := indexer.Search(query)
	if err != nil {
		return nil, err
	}

	feed := &torznab.ResultFeed{
		Info:  indexer.Info(),
		Items: items,
	}

	k, err := h.sharedKey()
	if err != nil {
		return nil, err
	}

	// rewrite links to use the server
	for idx, item := range feed.Items {
		t := &token{
			Site: item.Site,
			Link: item.Link,
		}

		te, err := t.Encode(k)
		if err != nil {
			return nil, err
		}
		baseURL.Path += fmt.Sprintf("/%s/%s.torrent", te, item.Title)
		feed.Items[idx].Link = baseURL.String()
	}

	return feed, err
}
