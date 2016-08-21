package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/cardigann/cardigann/config"
	"github.com/cardigann/cardigann/indexer"
	"github.com/cardigann/cardigann/server"
	"github.com/cardigann/cardigann/torznab"
)

var (
	Version string
)

func main() {
	os.Exit(run(os.Args[1:]...))
}

func run(args ...string) (exitCode int) {
	app := kingpin.New("cardigann",
		`A torznab proxy for torrent indexer sites`)

	app.Version(Version)
	app.Writer(os.Stdout)

	app.Terminate(func(code int) {
		exitCode = code
	})

	configureQueryCommand(app)
	configureDownloadCommand(app)
	configureServerCommand(app)
	configureTestDefinitionCommand(app)

	kingpin.MustParse(app.Parse(args))
	return
}

func lookupIndexer(key string) (torznab.Indexer, error) {
	conf, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	def, err := indexer.LoadDefinition(key)
	if err != nil {
		return nil, err
	}

	return indexer.NewRunner(def, conf)
}

func configureQueryCommand(app *kingpin.Application) {
	var key, format string
	var args []string

	cmd := app.Command("query", "Manually query an indexer using torznab commands")
	cmd.Alias("q")
	cmd.Flag("format", "Either json, xml or rss").
		Default("json").
		Short('f').
		EnumVar(&format, "xml", "json", "rss")

	cmd.Arg("key", "The indexer key").
		Required().
		StringVar(&key)

	cmd.Arg("args", "Arguments to use to query").
		StringsVar(&args)

	cmd.Action(func(c *kingpin.ParseContext) error {
		return queryCommand(key, format, args)
	})
}

func queryCommand(key, format string, args []string) error {
	indexer, err := lookupIndexer(key)
	if err != nil {
		return err
	}

	vals := url.Values{}
	for _, arg := range args {
		tokens := strings.SplitN(arg, "=", 2)
		if len(tokens) == 1 {
			vals.Set("q", tokens[0])
		} else {
			vals.Add(tokens[0], tokens[1])
		}
	}

	query, err := torznab.ParseQuery(vals)
	if err != nil {
		return fmt.Errorf("Parsing query failed: %s", err.Error())
	}

	feed, err := indexer.Search(query)
	if err != nil {
		return fmt.Errorf("Searching failed: %s", err.Error())
	}

	switch format {
	case "xml":
		x, err := xml.MarshalIndent(feed, "", "  ")
		if err != nil {
			return fmt.Errorf("Failed to marshal XML: %s", err.Error())
		}
		fmt.Printf("%s", x)

	case "json":
		j, err := json.MarshalIndent(feed, "", "  ")
		if err != nil {
			return fmt.Errorf("Failed to marshal JSON: %s", err.Error())
		}
		fmt.Printf("%s", j)
	}

	return nil
}

func configureDownloadCommand(app *kingpin.Application) {
	var key, url, file string

	cmd := app.Command("download", "Download a torrent from the tracker")
	cmd.Arg("key", "The indexer key").
		Required().
		StringVar(&key)

	cmd.Arg("url", "The url of the file to download").
		Required().
		StringVar(&url)

	cmd.Arg("file", "The filename to download to").
		Required().
		StringVar(&file)

	cmd.Action(func(c *kingpin.ParseContext) error {
		return downloadCommand(key, url, file)
	})
}

func downloadCommand(key, url, file string) error {
	indexer, err := lookupIndexer(key)
	if err != nil {
		return err
	}

	rc, _, err := indexer.Download(url)
	if err != nil {
		return fmt.Errorf("Downloading failed: %s", err.Error())
	}

	defer rc.Close()

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("Creating file failed: %s", err.Error())
	}

	n, err := io.Copy(f, rc)
	if err != nil {
		return fmt.Errorf("Creating file failed: %s", err.Error())
	}

	log.Printf("Downloaded %d bytes", n)
	return nil
}

func configureServerCommand(app *kingpin.Application) {
	var bindPort, bindAddr, password string
	var devMode bool

	cmd := app.Command("server", "Run the proxy (and web) server")
	cmd.Flag("port", "The port to listen on").
		OverrideDefaultFromEnvar("PORT").
		Default("5060").
		StringVar(&bindPort)

	cmd.Flag("addr", "The address to listen on").
		Default("0.0.0.0").
		StringVar(&bindAddr)

	cmd.Flag("passphrase", "Require a passphrase to view web interface").
		Short('p').
		Required().
		StringVar(&password)

	cmd.Flag("dev", "Run in local development mode").
		BoolVar(&devMode)

	cmd.Action(func(c *kingpin.ParseContext) error {
		return serverCommand(bindAddr, bindPort, password, devMode)
	})
}

func serverCommand(addr, port string, password string, devMode bool) error {
	conf, err := indexer.NewConfig()
	if err != nil {
		return err
	}

	listenOn := fmt.Sprintf("%s:%s", addr, port)
	log.Printf("Starting server on http://%s", listenOn)

	return http.ListenAndServe(listenOn, server.NewHandler(indexer.Registered, server.Params{
		DevMode:    devMode,
		Passphrase: password,
		Config:     conf,
	}))
}

func configureTestDefinitionCommand(app *kingpin.Application) {
	var f *os.File

	cmd := app.Command("test-definition", "Test a yaml indexer definition file")
	cmd.Alias("test")

	cmd.Arg("file", "The definition yaml file").
		Required().
		FileVar(&f)

	cmd.Action(func(c *kingpin.ParseContext) error {
		return testDefinitionCommand(f)
	})
}

func testDefinitionCommand(f *os.File) error {
	conf, err := indexer.NewConfig()
	if err != nil {
		return err
	}

	def, err := indexer.ParseDefinitionFile(f)
	if err != nil {
		return err
	}

	fmt.Println("Definition file parsing OK")

	runner := indexer.NewRunner(def, conf)

	err = runner.Login()
	if err != nil {
		return fmt.Errorf("Login failed: %s", err.Error())
	}

	fmt.Println("Login OK")

	err = runner.Test()
	if err != nil {
		return err
	}

	fmt.Println("Indexer test returned OK")
	return nil
}
