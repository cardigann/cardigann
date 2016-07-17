package server

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/cardigann/cardigann/indexer"
	"github.com/cardigann/cardigann/torznab"
	"github.com/julienschmidt/httprouter"

	// indexers
	_ "github.com/cardigann/cardigann/indexer/bithdtv"
)

type handler struct {
	Indexers indexer.ConstructorMap
	Params   Params
}

func (h *handler) BaseURL(r *http.Request) (*url.URL, error) {
	if h.Params.BaseURL != "" {
		return url.Parse(h.Params.BaseURL)
	}
	proto := "http"
	if r.TLS != nil {
		proto = "https"
	}
	return url.Parse(fmt.Sprintf("%s://%s", proto, r.Host))
}

func (h *handler) DownloadURL(r *http.Request, item torznab.ResultItem, siteKey string) (string, error) {
	ul, err := url.Parse(item.Link)
	if err != nil {
		return "", err
	}

	u, err := h.BaseURL(r)
	if err != nil {
		return "", err
	}

	slug := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s?%s", ul.Path, ul.RawQuery)))
	u.Path = fmt.Sprintf("/dl/%s/%s/%s.torrent", siteKey, slug, item.Title)

	return u.String(), nil
}

func (h *handler) IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Interactive UI not implemented yet\n")
}

func (h *handler) TorznabHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	indexer, err := indexer.Registered.New(ps.ByName("sitekey"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("GET %s", r.URL.String())
	t := r.URL.Query().Get("t")

	switch t {
	case "caps":
		indexer.Capabilities().ServeHTTP(w, r)

	case "search", "tvsearch", "tv-search":
		feed, err := h.Search(r, indexer, ps.ByName("sitekey"))
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

func (h *handler) DownloadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	indexer, err := indexer.Registered.New(ps.ByName("sitekey"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	log.Printf("GET %s", r.URL.String())

	slug, err := base64.StdEncoding.DecodeString(ps.ByName("filekey"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := url.Parse(string(slug))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rc, err := indexer.Download(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	defer rc.Close()
	io.Copy(w, rc)
}

func (h *handler) Search(r *http.Request, indexer torznab.Indexer, siteKey string) (*torznab.ResultFeed, error) {
	query, err := torznab.ParseQuery(r.URL.Query())
	if err != nil {
		return nil, err
	}

	log.Printf("Query: %#v", query)
	feed, err := indexer.(torznab.Indexer).Search(query)
	if err != nil {
		return nil, err
	}

	// rewrite links to use the server
	for idx, item := range feed.Items {
		dl, err := h.DownloadURL(r, item, siteKey)
		if err != nil {
			return nil, err
		}
		feed.Items[idx].Link = dl
	}

	return feed, err
}

type Params struct {
	BaseURL string
}

func ListenAndServe(listenAddr string, cm indexer.ConstructorMap, p Params) error {
	h := handler{Indexers: cm, Params: p}

	router := httprouter.New()
	router.GET("/", h.IndexHandler)
	router.GET("/torznab/:sitekey/api", h.TorznabHandler)
	router.GET("/dl/:sitekey/:filekey/:filename", h.DownloadHandler)

	log.Printf("Listening on %s", listenAddr)
	return http.ListenAndServe(listenAddr, router)
}
