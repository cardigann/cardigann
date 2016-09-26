//go:generate esc -o static.go -prefix ../web/build -pkg server ../web/build
package server

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cardigann/cardigann/config"
	"github.com/cardigann/cardigann/indexer"
	"github.com/cardigann/cardigann/logger"
	"github.com/cardigann/cardigann/torznab"
	"github.com/gorilla/mux"
)

const (
	buildDir = "/web/build"
)

var (
	log              = logger.Logger
	apiRoutePrefixes = []string{
		"/torznab/",
		"/download/",
		"/xhr/",
		"/debug/",
	}
)

type Params struct {
	BaseURL    string
	APIKey     []byte
	Passphrase string
	Config     config.Config
	Version    string
}

type handler struct {
	http.Handler
	Params      Params
	FileHandler http.Handler
	indexers    map[string]torznab.Indexer
}

func NewHandler(p Params) (http.Handler, error) {
	h := &handler{
		Params:      p,
		FileHandler: http.FileServer(FS(false)),
		indexers:    map[string]torznab.Indexer{},
	}

	router := mux.NewRouter()

	// torznab routes
	router.HandleFunc("/torznab/{indexer}", h.torznabHandler).Methods("GET")
	router.HandleFunc("/torznab/{indexer}/api", h.torznabHandler).Methods("GET")
	router.HandleFunc("/download/{indexer}/{token}/{filename}", h.downloadHandler).Methods("HEAD")
	router.HandleFunc("/download/{indexer}/{token}/{filename}", h.downloadHandler).Methods("GET")
	router.HandleFunc("/download/{token}/{filename}", h.downloadHandler).Methods("HEAD")
	router.HandleFunc("/download/{token}/{filename}", h.downloadHandler).Methods("GET")

	// xhr routes for the webapp
	router.HandleFunc("/xhr/indexers/{indexer}/test", h.getIndexerTestHandler).Methods("GET")
	router.HandleFunc("/xhr/indexers/{indexer}/config", h.getIndexersConfigHandler).Methods("GET")
	router.HandleFunc("/xhr/indexers/{indexer}/config", h.patchIndexersConfigHandler).Methods("PATCH")
	router.HandleFunc("/xhr/indexers", h.getIndexersHandler).Methods("GET")
	router.HandleFunc("/xhr/indexers", h.patchIndexersHandler).Methods("PATCH")
	router.HandleFunc("/xhr/auth", h.getAuthHandler).Methods("GET")
	router.HandleFunc("/xhr/auth", h.postAuthHandler).Methods("POST")
	router.HandleFunc("/xhr/version", h.getVersionHandler).Methods("GET")

	h.Handler = router
	return h, h.initialize()
}

func (h *handler) initialize() error {
	if h.Params.Passphrase == "" {
		apiKey, hasApiKey, err := h.Params.Config.Get("global", "apikey")
		if err != nil {
			return err
		}
		if !hasApiKey {
			k, err := h.sharedKey()
			if err != nil {
				return err
			}
			h.Params.APIKey = k
			return h.Params.Config.Set("global", "apikey", fmt.Sprintf("%x", k))
		}
		k, err := hex.DecodeString(apiKey)
		if err != nil {
			return err
		}
		h.Params.APIKey = k
	}
	return nil
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

func (h *handler) createIndexer(key string) (torznab.Indexer, error) {
	def, err := indexer.DefaultDefinitionLoader.Load(key)
	if err != nil {
		log.WithError(err).Warnf("Failed to load definition for %q", key)
		return nil, err
	}

	log.WithFields(logrus.Fields{"indexer": key}).Debugf("Loaded indexer")
	indexer, err := indexer.NewRunner(def, indexer.RunnerOpts{
		Config: h.Params.Config,
	}), nil
	if err != nil {
		return nil, err
	}

	return indexer, nil
}

func (h *handler) lookupIndexer(key string) (torznab.Indexer, error) {
	if key == "aggregate" {
		return h.createAggregate()
	}
	if _, ok := h.indexers[key]; !ok {
		indexer, err := h.createIndexer(key)
		if err != nil {
			return nil, err
		}
		h.indexers[key] = indexer
	}

	return h.indexers[key], nil
}

func (h *handler) createAggregate() (torznab.Indexer, error) {
	keys, err := indexer.DefaultDefinitionLoader.List()
	if err != nil {
		return nil, err
	}

	agg := indexer.Aggregate{}
	for _, key := range keys {
		if config.IsSectionEnabled(key, h.Params.Config) {
			indexer, err := h.lookupIndexer(key)
			if err != nil {
				return nil, err
			}
			agg = append(agg, indexer)
		}
	}

	return agg, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Cache-Control, Content-Type, Content-Length, Accept-Encoding, Authorization, Last-Event-ID")
	}
	if r.Method == "OPTIONS" {
		return
	}

	log.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.RequestURI(),
		"remote": r.RemoteAddr,
	}).Debugf("%s %s", r.Method, r.URL.RequestURI())

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

	indexer, err := h.lookupIndexer(indexerID)
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
		switch r.URL.Query().Get("format") {
		case "", "xml":
			x, err := xml.MarshalIndent(feed, "", "  ")
			if err != nil {
				torznab.Error(w, err.Error(), torznab.ErrUnknownError)
				return
			}
			w.Header().Set("Content-Type", "application/rss+xml")
			w.Write(x)
		case "json":
			jsonOutput(w, feed)
		}

	default:
		torznab.Error(w, "Unknown type parameter", torznab.ErrIncorrectParameter)
	}
}

func (h *handler) downloadHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	token := params["token"]
	filename := params["filename"]

	log.WithFields(logrus.Fields{"filename": filename}).Debugf("Processing download via handler")

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

	indexer, err := h.lookupIndexer(t.Site)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	rc, _, err := indexer.Download(t.Link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/x-bittorrent")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Transfer-Encoding", "binary")

	defer rc.Close()
	io.Copy(w, rc)
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
			log.Debugf("Error encoding token: %v", err)
			return nil, err
		}

		log.Debugf("Generated signed token %q", te)
		feed.Items[idx].Link = fmt.Sprintf("%s/%s/%s.torrent", baseURL.String(), te, item.Title)
	}

	return feed, err
}
