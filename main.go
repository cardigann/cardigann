package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/cardigann/cardigann/indexer"
	"github.com/cardigann/cardigann/torznab"
)

var (
	app = kingpin.New("cardigann", "A proxy for private trackers")

	query       = app.Command("query", "Query an indexer")
	queryFormat = query.Flag("format", "Either json, xml or rss").Default("json").Enum("xml", "json", "rss")
	queryKey    = query.Arg("key", "The indexer key").Required().String()
	queryArgs   = query.Arg("args", "Arguments to use to query").Strings()
)

func queryCommand() {
	indexer, err := indexer.Get(*queryKey)
	if err != nil {
		kingpin.Fatalf(err.Error())
	}

	query := make(torznab.Query)
	for _, arg := range *queryArgs {
		tokens := strings.SplitN(arg, "=", 2)
		query[tokens[0]] = tokens[1]
	}

	items, err := indexer.(torznab.Indexer).Search(query)
	if err != nil {
		kingpin.Fatalf("Searching failed: %s", err.Error())
	}

	switch *queryFormat {
	case "xml":
		feed := torznab.ResultFeed{
			Items: items,
		}

		x, err := xml.MarshalIndent(feed, "", "  ")
		if err != nil {
			kingpin.Fatalf("Failed to marshal XML: %s", err.Error())
		}
		fmt.Printf("%s", x)

	case "json":
		j, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			kingpin.Fatalf("Failed to marshal JSON: %s", err.Error())
		}
		fmt.Printf("%s", j)
	}
}

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case query.FullCommand():
		queryCommand()
	}
}
