package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/cardigann/cardigann/indexer"
	"github.com/cardigann/cardigann/torznab"

	// indexers
	_ "github.com/cardigann/cardigann/indexer/bithdtv"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Interactive UI not implemented yet\n")
}

func torznabHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	indexer, err := indexer.Get(ps.ByName("sitekey"))
	if err != nil {
		log.Fatal(err)
	}

	t := r.URL.Query().Get("t")

	switch t {
	case "caps":
		indexer.Capabilities().ServeHTTP(w, r)
	case "search", "tv-search":
		if !indexer.Capabilities().HasSearchMode(t) {
			http.Error(w, fmt.Sprintf("Indexer %q doesn't support t=%q", ps.ByName("sitekey"), t), http.StatusNotFound)
			return
		}
		query, err := torznab.ParseQuery(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "<h1>Query</h1><pre><code>%#v</code></pre>", query)
		res, err := indexer.(torznab.Indexer).Search(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		fmt.Fprintf(w, "<h1>Result</h1><pre><code>%#v</code></pre>", res)
	default:
		http.Error(w, fmt.Sprintf("Type %q not implemented", t), http.StatusNotFound)
	}
}

func main() {
	log.Println("Listening on :3000")

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/api/:sitekey/torznab", torznabHandler)

	log.Fatal(http.ListenAndServe(":3000", router))
}
