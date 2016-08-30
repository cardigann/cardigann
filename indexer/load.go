package indexer

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/cardigann/cardigann/config"
)

var ErrUnknownIndexer = errors.New("Unknown indexer")

func findDefinitions() (map[string]string, error) {
	dirs, err := config.Dirs()
	if err != nil {
		return nil, err
	}

	results := map[string]string{}
	for _, dirpath := range dirs {
		if dir, err := os.Open(path.Join(dirpath, "definitions")); err == nil {
			files, err := dir.Readdirnames(-1)
			if err != nil {
				continue
			}
			for _, basename := range files {
				if strings.HasSuffix(basename, ".yml") {
					results[strings.TrimSuffix(basename, ".yml")] = path.Join(dir.Name(), basename)
				}
			}
		}
	}

	return results, nil
}

func ListDefinitions() ([]string, error) {
	keys, err := findDefinitions()
	if err != nil {
		return nil, err
	}

	results := []string{}

	for k := range keys {
		results = append(results, k)
	}

	return results, nil
}

func LoadDefinition(key string) (*IndexerDefinition, error) {
	defs, err := findDefinitions()
	if err != nil {
		return nil, err
	}

	fileName, ok := defs[key]
	if !ok {
		return nil, ErrUnknownIndexer
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return ParseDefinition(data)
}
