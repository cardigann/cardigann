package indexer

import (
	"fmt"

	"github.com/cardigann/cardigann/torznab"
	"github.com/vaughan0/go-ini"
)

var config Config

type Config interface {
	Get(section string, key string) (string, bool)
}

var TorznabIndexers = map[string]func(c Config) (torznab.Indexer, error){}

// Indexer creates a new indexer with configuration
func Get(key string) (torznab.Indexer, error) {
	if config == nil {
		var err error
		config, err = ini.LoadFile("config.ini")

		if err != nil {
			return nil, err
		}
	}

	indexerFunc, ok := TorznabIndexers[key]
	if !ok {
		return nil, fmt.Errorf("Indexer %s doesn't exist", key)
	}

	indexer, err := indexerFunc(config)
	if err != nil {
		return nil, err
	}

	return indexer, nil
}
