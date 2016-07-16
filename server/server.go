package server

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

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

func (h *handler) BaseURL(r *http.Request) string {
	if h.Params.BaseURL != "" {
		return h.Params.BaseURL
	}
	proto := "http"
	if r.TLS != nil {
		proto = "https"
	}
	return fmt.Sprintf("%s://%s%s", proto, r.Host, r.URL.Path)
}

func (h *handler) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	case "search", "tvsearch":
		feed, err := h.Search(r, indexer)
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

func (h *handler) Search(r *http.Request, indexer torznab.Indexer) (*torznab.ResultFeed, error) {
	query, err := torznab.ParseQuery(r.URL.Query())
	if err != nil {
		return nil, err
	}

	log.Printf("Query: %#v", query)
	feed, err := indexer.(torznab.Indexer).Search(query)
	if err != nil {
		return nil, err
	}

	// feed.Self = h.BaseURL(r)
	return feed, err
}

type Params struct {
	BaseURL string
}

func ListenAndServe(listenAddr string, cm indexer.ConstructorMap, p Params) error {
	h := handler{Indexers: cm, Params: p}

	router := httprouter.New()
	router.GET("/", h.Index)
	router.GET("/api/:sitekey/torznab", h.TorznabHandler)
	router.GET("/api/:sitekey/torznab/api", h.TorznabHandler)

	log.Printf("Listening on %s", listenAddr)
	return http.ListenAndServe(listenAddr, router)
}
