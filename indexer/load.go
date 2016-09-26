package indexer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/cardigann/cardigann/config"
)

var (
	ErrUnknownIndexer       = errors.New("Unknown indexer")
	DefaultDefinitionLoader DefinitionLoader
)

type DefinitionLoader interface {
	List() ([]string, error)
	Load(key string) (*IndexerDefinition, error)
}

func init() {
	DefaultDefinitionLoader = &multiLoader{
		newFsLoader(),
		escLoader{Dir(false, "")},
	}
}

type fsLoader struct {
	dirs []string
}

func newFsLoader() DefinitionLoader {
	return &fsLoader{config.GetDefinitionDirs()}
}

func (fs *fsLoader) walkDirectories() (map[string]string, error) {
	defs := map[string]string{}

	for _, dirpath := range fs.dirs {
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
				defs[strings.TrimSuffix(basename, ".yml")] = path.Join(dir.Name(), basename)
			}
		}
	}

	return defs, nil
}

func (fs *fsLoader) List() ([]string, error) {
	defs, err := fs.walkDirectories()
	if err != nil {
		return nil, err
	}

	results := []string{}

	for k := range defs {
		results = append(results, k)
	}

	return results, nil
}

func (fs *fsLoader) Load(key string) (*IndexerDefinition, error) {
	defs, err := fs.walkDirectories()
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

type multiLoader []DefinitionLoader

func (ml multiLoader) List() ([]string, error) {
	results := []string{}

	for _, loader := range ml {
		result, err := loader.List()
		if err != nil {
			return nil, err
		}
		results = append(results, result...)
	}

	return results, nil
}

func (ml multiLoader) Load(key string) (*IndexerDefinition, error) {
	for _, loader := range ml {
		def, err := loader.Load(key)
		if err == nil {
			return def, nil
		}
	}
	return nil, ErrUnknownIndexer
}

type escLoader struct {
	http.FileSystem
}

func (el escLoader) List() ([]string, error) {
	results := []string{}

	for key := range _escData {
		results = append(results, key)
	}

	return results, nil
}

func (el escLoader) Load(key string) (*IndexerDefinition, error) {
	f, err := el.Open(fmt.Sprintf("/definitions/%s.yml", key))
	if os.IsNotExist(err) {
		return nil, ErrUnknownIndexer
	} else if err != nil {
		return nil, err
	}

	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return ParseDefinition(data)
}
