package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "net/http/pprof"

	"github.com/Sirupsen/logrus"
	"github.com/cardigann/cardigann/config"
	"github.com/cardigann/cardigann/indexer"
	"github.com/cardigann/cardigann/logger"
	"github.com/cardigann/cardigann/server"
	"github.com/cardigann/cardigann/torznab"
	"github.com/kardianos/service"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Version string
	log     = logger.Logger
)

func main() {
	os.Exit(run(os.Args[1:]...))
}

func run(args ...string) (exitCode int) {
	app := kingpin.New("cardigann",
		`A torznab proxy for torrent indexer sites`)

	app.Writer(os.Stdout)
	app.DefaultEnvars()

	app.Terminate(func(code int) {
		exitCode = code
	})

	if err := configureServerCommand(app); err != nil {
		log.Error(err)
		return 1
	}

	configureQueryCommand(app)
	configureDownloadCommand(app)
	configureTestDefinitionCommand(app)
	configureServiceCommand(app)

	app.Command("version", "Print the application version").Action(func(c *kingpin.ParseContext) error {
		fmt.Print(Version)
		return nil
	})

	kingpin.MustParse(app.Parse(args))
	return
}

func newConfig() (config.Config, error) {
	f, err := config.GetConfigPath()
	if err != nil {
		return nil, err
	}

	log.WithField("path", f).Debug("Reading config")
	return config.NewJSONConfig(f)
}

func lookupRunner(key string, opts indexer.RunnerOpts) (torznab.Indexer, error) {
	if key == "aggregate" {
		return lookupAggregate(opts)
	}

	def, err := indexer.DefaultDefinitionLoader.Load(key)
	if err != nil {
		return nil, err
	}

	return indexer.NewRunner(def, opts), nil
}

func lookupAggregate(opts indexer.RunnerOpts) (torznab.Indexer, error) {
	keys, err := indexer.DefaultDefinitionLoader.List()
	if err != nil {
		return nil, err
	}

	agg := indexer.Aggregate{}
	for _, key := range keys {
		if config.IsSectionEnabled(key, opts.Config) {
			def, err := indexer.DefaultDefinitionLoader.Load(key)
			if err != nil {
				return nil, err
			}

			agg = append(agg, indexer.NewRunner(def, opts))
		}
	}

	return agg, nil
}

var globals struct {
	Debug bool
}

func configureGlobalFlags(cmd *kingpin.CmdClause) {
	cmd.Flag("debug", "Print out debug logging").BoolVar(&globals.Debug)
}

func applyGlobalFlags() {
	if globals.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}
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

	configureGlobalFlags(cmd)

	cmd.Action(func(c *kingpin.ParseContext) error {
		applyGlobalFlags()
		return queryCommand(key, format, args)
	})
}

func queryCommand(key, format string, args []string) error {
	conf, err := newConfig()
	if err != nil {
		return err
	}

	indexer, err := lookupRunner(key, indexer.RunnerOpts{
		Config: conf,
	})
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

	configureGlobalFlags(cmd)

	cmd.Action(func(c *kingpin.ParseContext) error {
		applyGlobalFlags()
		return downloadCommand(key, url, file)
	})
}

func downloadCommand(key, url, file string) error {
	conf, err := newConfig()
	if err != nil {
		return err
	}

	indexer, err := lookupRunner(key, indexer.RunnerOpts{
		Config: conf,
	})
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

	log.WithFields(logrus.Fields{"bytes": n}).Info("Downloading file")
	return nil
}

func configureServerCommand(app *kingpin.Application) error {
	var bindPort, bindAddr, password string

	conf, err := newConfig()
	if err != nil {
		return err
	}

	defaultBind, err := config.GetGlobalConfig("bind", "0.0.0.0", conf)
	if err != nil {
		return err
	}

	defaultPort, err := config.GetGlobalConfig("port", "5060", conf)
	if err != nil {
		return err
	}

	cmd := app.Command("server", "Run the proxy (and web) server")
	cmd.Flag("port", "The port to listen on").
		OverrideDefaultFromEnvar("PORT").
		Default(defaultPort).
		StringVar(&bindPort)

	cmd.Flag("bind", "The address to bind to").
		Default(defaultBind).
		StringVar(&bindAddr)

	cmd.Flag("passphrase", "Require a passphrase to view web interface").
		Short('p').
		StringVar(&password)

	configureGlobalFlags(cmd)
	cmd.Action(func(c *kingpin.ParseContext) error {
		applyGlobalFlags()
		return serverCommand(bindAddr, bindPort, password)
	})

	return nil
}

func serverCommand(addr, port string, password string) error {
	if globals.Debug {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	v := Version
	if v == "" {
		v = "dev"
	}

	log.Infof("Cardigann %s", v)

	conf, err := newConfig()
	if err != nil {
		return err
	}

	for _, dir := range config.GetDefinitionDirs() {
		log.WithField("dir", dir).Debug("Adding dir to definition load path")
	}

	listenOn := fmt.Sprintf("%s:%s", addr, port)
	log.Infof("Listening on %s", listenOn)

	h, err := server.NewHandler(server.Params{
		Passphrase: password,
		Config:     conf,
		Version:    Version,
	})
	if err != nil {
		return err
	}

	return http.ListenAndServe(listenOn, h)
}

func configureTestDefinitionCommand(app *kingpin.Application) {
	var f *os.File
	var cachePages bool

	cmd := app.Command("test-definition", "Test a yaml indexer definition file")
	cmd.Alias("test")

	cmd.Flag("cachepages", "Whether to store the output of browser actions for debugging").
		BoolVar(&cachePages)

	cmd.Arg("file", "The definition yaml file").
		Required().
		FileVar(&f)

	configureGlobalFlags(cmd)
	cmd.Action(func(c *kingpin.ParseContext) error {
		applyGlobalFlags()
		return testDefinitionCommand(f, cachePages)
	})
}

func testDefinitionCommand(f *os.File, cachePages bool) error {
	conf, err := newConfig()
	if err != nil {
		return err
	}

	def, err := indexer.ParseDefinitionFile(f)
	if err != nil {
		return err
	}

	fmt.Println("Definition file parsing OK")

	runner := indexer.NewRunner(def, indexer.RunnerOpts{
		Config:     conf,
		CachePages: cachePages,
	})
	tester := indexer.Tester{Runner: runner, Opts: indexer.TesterOpts{
		Download: true,
	}}

	err = tester.Test()
	if err != nil {
		return fmt.Errorf("Test failed: %s", err.Error())
	}

	fmt.Println("Indexer test returned OK")
	return nil
}

func configureServiceCommand(app *kingpin.Application) {
	var action string
	var userService bool
	var possibleActions = append(service.ControlAction[:], "run")

	cmd := app.Command("service", "Control the cardigann service")

	cmd.Flag("user", "Whether to use a user service rather than a system one").
		BoolVar(&userService)

	cmd.Arg("action", "One of "+strings.Join(possibleActions, ", ")).
		Required().
		EnumVar(&action, possibleActions...)

	configureGlobalFlags(cmd)
	cmd.Action(func(c *kingpin.ParseContext) error {
		log.Debugf("Running service action %s on platform %v.", action, service.Platform())

		conf, err := newConfig()
		if err != nil {
			return err
		}

		prg, err := newProgram(programOpts{
			UserService: userService,
			Config:      conf,
		})
		if err != nil {
			return err
		}

		if action != "run" {
			return service.Control(prg.service, action)
		}

		return runServiceCommand(prg)
	})
}

func runServiceCommand(prg *program) error {
	var err error
	errs := make(chan error)
	prg.logger, err = prg.service.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	logger.SetOutput(ioutil.Discard)
	logger.AddHook(&serviceLogHook{prg.logger})
	logger.SetFormatter(&serviceLogFormatter{})

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Error(err)
			}
		}
	}()

	err = prg.service.Run()
	if err != nil {
		prg.logger.Error(err)
	}

	return nil
}
