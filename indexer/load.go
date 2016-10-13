package indexer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/cardigann/cardigann/config"
)

var (
	ErrUnknownIndexer       = errors.New("Unknown indexer")
	DefaultDefinitionLoader DefinitionLoader
)

func ListBuiltins() ([]string, error) {
	l := escLoader{Dir(false, "")}
	return l.List()
}

func LoadEnabledDefinitions(conf config.Config) ([]*IndexerDefinition, error) {
	keys, err := DefaultDefinitionLoader.List()
	if err != nil {
		return nil, err
	}
	defs := []*IndexerDefinition{}
	for _, key := range keys {
		if config.IsSectionEnabled(key, conf) {
			def, err := DefaultDefinitionLoader.Load(key)
			if err != nil {
				return nil, err
			}
			defs = append(defs, def)
		}
	}
	return defs, nil
}

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

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	def, err := ParseDefinitionFile(f)
	if err != nil {
		return def, err
	}

	def.stats.Source = "file:" + fileName
	return def, err
}

type multiLoader []DefinitionLoader

func (ml multiLoader) List() ([]string, error) {
	allResults := map[string]struct{}{}

	for _, loader := range ml {
		result, err := loader.List()
		if err != nil {
			return nil, err
		}
		for _, val := range result {
			allResults[val] = struct{}{}
		}
	}

	results := []string{}

	for key := range allResults {
		results = append(results, key)
	}

	sort.Sort(sort.StringSlice(results))
	return results, nil
}

func (ml multiLoader) Load(key string) (*IndexerDefinition, error) {
	var def *IndexerDefinition

	for _, loader := range ml {
		loaded, err := loader.Load(key)
		if err != nil {
			continue
		}
		if def == nil || loaded.Stats().ModTime.After(def.Stats().ModTime) {
			def = loaded
		}
	}

	if def == nil {
		return nil, ErrUnknownIndexer
	}

	return def, nil
}

var escFilenameRegex = regexp.MustCompile(`^/definitions/(.+?)\.yml$`)

type escLoader struct {
	http.FileSystem
}

func (el escLoader) List() ([]string, error) {
	results := []string{}

	for filename := range _escData {
		if matches := escFilenameRegex.FindStringSubmatch(filename); matches != nil {
			results = append(results, matches[1])
		}
	}

	return results, nil
}

func (el escLoader) Load(key string) (*IndexerDefinition, error) {
	fname := fmt.Sprintf("/definitions/%s.yml", key)
	f, err := el.Open(fname)
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

	def, err := ParseDefinition(data)
	if err != nil {
		return def, err
	}

	fi, err := f.Stat()
	if err != nil {
		return def, err
	}

	def.stats.ModTime = fi.ModTime()
	def.stats.Source = "builtin:" + fname
	return def, nil
}
