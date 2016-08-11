package indexer

import (
	"fmt"

	"github.com/cardigann/cardigann/torznab"
)

var (
	Registered = make(ConstructorMap)
)

type Constructor func(c Config) (torznab.Indexer, error)
type ConstructorMap map[string]Constructor

// New creates a new torznab indexer with config
func (c ConstructorMap) New(key string, config Config) (torznab.Indexer, error) {
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
