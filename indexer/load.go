package indexer

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/cardigann/cardigann/config"
	"github.com/cardigann/cardigann/logger"
)

var definitions = map[string]string{}
var ErrUnknownIndexer = errors.New("Unknown indexer")

func findDefinitions() (map[string]string, error) {
	dirs, err := config.GetDefinitionDirs()
	if err != nil {
		return nil, err
	}

	for _, dirpath := range dirs {
		dir, err := os.Open(dirpath)
		if os.IsNotExist(err) {
			continue
		}
		files, err := dir.Readdirnames(-1)
		if err != nil {
			continue
		}
		for _, basename := range files {
			if strings.HasSuffix(basename, ".yml") {
				key := strings.TrimSuffix(basename, ".yml")
				f := path.Join(dir.Name(), basename)
				if existing := definitions[key]; existing != f {
					logger.Logger.WithField("path", f).Debug("Found definition")
					definitions[key] = f
				}
			}
		}
	}

	return definitions, nil
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
