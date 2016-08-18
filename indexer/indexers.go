package indexer

import (
	"errors"

	"github.com/cardigann/cardigann/torznab"
)

var (
	Registered = make(ConstructorMap)
)

var ErrUnknownIndexer = errors.New("Unknown indexer")

type Constructor func(c Config) (torznab.Indexer, error)
type ConstructorMap map[string]Constructor

// New creates a new torznab indexer with config
func (c ConstructorMap) New(key string, config Config) (torznab.Indexer, error) {
	indexerFunc, ok := c[key]
	if !ok {
		return nil, ErrUnknownIndexer
	}

	indexer, err := indexerFunc(config)
	if err != nil {
		return nil, err
	}

	return indexer, nil
}
