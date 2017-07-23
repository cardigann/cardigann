//go:generate esc -o static.go -prefix ../web/build -pkg server ../web/build
package server

import (
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/cardigann/cardigann/config"
	"github.com/cardigann/cardigann/indexer"
	"github.com/cardigann/cardigann/logger"
	"github.com/cardigann/cardigann/torrentpotato"
	"github.com/cardigann/cardigann/torznab"
	"github.com/gorilla/mux"
)

const (
	buildDir = "/web/build"
)

var (
	log = logger.Logger
)

type Params struct {
	BaseURL    string
	PathPrefix string
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
		Params: p,
		FileHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Loading %s from fs", r.URL.RequestURI())
			http.FileServer(FS(false)).ServeHTTP(w, r)
		}),
		indexers: map[string]torznab.Indexer{},
	}

	router := mux.NewRouter()

	// apply the path prefix
	if h.Params.PathPrefix != "/" && h.Params.PathPrefix != "" {
		router = router.PathPrefix(h.Params.PathPrefix).Subrouter()
		h.FileHandler = http.StripPrefix(h.Params.PathPrefix, h.FileHandler)
	}

	// torznab routes
	router.HandleFunc("/torznab/{indexer}", h.torznabHandler).Methods("GET")
	router.HandleFunc("/torznab/{indexer}/api", h.torznabHandler).Methods("GET")

	// torrentpotato routes
	router.HandleFunc("/torrentpotato/{indexer}", h.torrentPotatoHandler).Methods("GET")

	// download routes
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

	// anything else
	router.PathPrefix("/").Handler(h.FileHandler)
	router.PathPrefix("/static").Handler(h.FileHandler)

	h.Handler = router
	return h, h.initialize()
}

func (h *handler) initialize() error {
	if h.Params.Passphrase == "" {
		pass, hasPassphrase, _ := h.Params.Config.Get("global", "passphrase")
		if hasPassphrase {
			h.Params.Passphrase = pass
			return nil
		}
		apiKey, hasApiKey, _ := h.Params.Config.Get("global", "apikey")
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

	// Walk routes for debugging
	err := h.Handler.(*mux.Router).Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		log.Debugf("Responds to %s", path)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *handler) baseURL(r *http.Request, appendPath string) (*url.URL, error) {
	proto := "http"
	if r.TLS != nil {
		proto = "https"
	}
	return url.Parse(fmt.Sprintf("%s://%s%s", proto, r.Host,
		path.Join(h.Params.PathPrefix, appendPath)))
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

	h.Handler.ServeHTTP(w, r)
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
		feed, err := h.torznabSearch(r, indexer, indexerID)
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

func (h *handler) torrentPotatoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	indexerID := params["indexer"]

	apiKey := r.URL.Query().Get("passkey")
	if !h.checkAPIKey(apiKey) {
		torrentpotato.Error(w, errors.New("Invalid passkey"))
		return
	}

	indexer, err := h.lookupIndexer(indexerID)
	if err != nil {
		torrentpotato.Error(w, err)
		return
	}

	query := torznab.Query{
		Type: "movie",
		Categories: []int{
			torznab.CategoryMovies.ID,
			torznab.CategoryMovies_SD.ID,
			torznab.CategoryMovies_HD.ID,
			torznab.CategoryMovies_Foreign.ID,
		},
	}

	qs := r.URL.Query()

	if search := qs.Get("search"); search != "" {
		query.Q = search
	}

	if imdbid := qs.Get("imdbid"); imdbid != "" {
		query.IMDBID = imdbid
	}

	items, err := indexer.Search(query)
	if err != nil {
		torrentpotato.Error(w, err)
		return
	}

	rewritten, err := h.rewriteLinks(r, items)
	if err != nil {
		torrentpotato.Error(w, err)
		return
	}

	torrentpotato.Output(w, rewritten)
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

func (h *handler) torznabSearch(r *http.Request, indexer torznab.Indexer, siteKey string) (*torznab.ResultFeed, error) {
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

	rewritten, err := h.rewriteLinks(r, items)
	if err != nil {
		return nil, err
	}

	feed.Items = rewritten
	return feed, err
}

func (h *handler) rewriteLinks(r *http.Request, items []torznab.ResultItem) ([]torznab.ResultItem, error) {
	baseURL, err := h.baseURL(r, "/download")
	if err != nil {
		return nil, err
	}

	k, err := h.sharedKey()
	if err != nil {
		return nil, err
	}

	// rewrite non-magnet links to use the server
	for idx, item := range items {
		if strings.HasPrefix(item.Link, "magnet:") {
			continue
		}

		t := &token{
			Site: item.Site,
			Link: item.Link,
		}

		te, err := t.Encode(k)
		if err != nil {
			log.Debugf("Error encoding token: %v", err)
			return nil, err
		}

		filename := strings.Replace(item.Title, "/", "-", -1)
		items[idx].Link = fmt.Sprintf("%s/%s/%s.torrent", baseURL.String(), te, url.QueryEscape(filename))
	}

	return items, nil
}
