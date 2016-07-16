package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/cardigann/cardigann/indexer"
	cserver "github.com/cardigann/cardigann/server"
	"github.com/cardigann/cardigann/torznab"

	// indexer drivers to load
	_ "github.com/cardigann/cardigann/indexer/bithdtv"
)

var (
	app = kingpin.New("cardigann", "A proxy for private trackers")

	query       = app.Command("query", "Query an indexer")
	queryFormat = query.Flag("format", "Either json, xml or rss").Default("json").Enum("xml", "json", "rss")
	queryKey    = query.Arg("key", "The indexer key").Required().String()
	queryArgs   = query.Arg("args", "Arguments to use to query").Strings()

	server     = app.Command("server", "Run the proxy server")
	serverAddr = server.Flag("addr", "The host and port to bind to").Default(":3000").String()
)

func queryCommand() {
	indexer, err := indexer.Registered.New(*queryKey)
	if err != nil {
		kingpin.Fatalf(err.Error())
	}

	query := make(torznab.Query)
	for _, arg := range *queryArgs {
		tokens := strings.SplitN(arg, "=", 2)
		query[tokens[0]] = tokens[1]
	}

	feed, err := indexer.(torznab.Indexer).Search(query)
	if err != nil {
		kingpin.Fatalf("Searching failed: %s", err.Error())
	}

	switch *queryFormat {
	case "xml":
		x, err := xml.MarshalIndent(feed, "", "  ")
		if err != nil {
			kingpin.Fatalf("Failed to marshal XML: %s", err.Error())
		}
		fmt.Printf("%s", x)

	case "json":
		j, err := json.MarshalIndent(feed, "", "  ")
		if err != nil {
			kingpin.Fatalf("Failed to marshal JSON: %s", err.Error())
		}
		fmt.Printf("%s", j)
	}
}

func serverCommand() {
	log.Fatal(cserver.ListenAndServe(*serverAddr, indexer.Registered, cserver.Params{}))
}

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case query.FullCommand():
		queryCommand()
	case server.FullCommand():
		serverCommand()
	}
}
