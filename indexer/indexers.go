package indexer

import (
	"fmt"

	"github.com/cardigann/cardigann/torznab"
	"github.com/vaughan0/go-ini"
)

var (
	Registered = make(ConstructorMap)
)

type Config interface {
	Get(section string, key string) (string, bool)
}

type Constructor func(c Config) (torznab.Indexer, error)
type ConstructorMap map[string]Constructor

// New creates a new torznab indexer with config loaded from config.ini
func (c ConstructorMap) New(key string) (torznab.Indexer, error) {
	config, err := ini.LoadFile("config.ini")
	if err != nil {
		return nil, err
	}

	indexerFunc, ok := c[key]
	if !ok {
		return nil, fmt.Errorf("Indexer %s doesn't exist", key)
	}

	indexer, err := indexerFunc(config)
	if err != nil {
		return nil, err
	}

	return indexer, nil
}
