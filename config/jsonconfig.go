package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/cardigann/cardigann/logger"
)

var (
	log = logger.Logger
)

const (
	configFileName = "config.json"
)

type jsonConfig struct {
	dirs       []string
	defaultDir string
}

type boolOrString string

func (b *boolOrString) UnmarshalJSON(data []byte) error {
	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err == nil {
		*b = boolOrString(strconv.FormatBool(boolVal))
		return nil
	}

	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		*b = boolOrString(strVal)
		return nil
	}

	return errors.New("Failed to unmarshal boolOrString")
}

type jsonConfigMap map[string]map[string]boolOrString

func NewJSONConfig() (Config, error) {
	dirs, err := Dirs()
	if err != nil {
		return nil, err
	}

	defaultDir, err := DefaultDir()
	if err != nil {
		return nil, err
	}

	return &jsonConfig{dirs, defaultDir}, nil
}

func (jc *jsonConfig) load() (jsonConfigMap, error) {
	config := jsonConfigMap{}

	path, err := Find(configFileName, jc.dirs)
	if err == nil {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, err
		}
	}
	return config, nil
}

func (jc *jsonConfig) save(c jsonConfigMap) error {
	path, err := Find(configFileName, jc.dirs)
	if err != nil {
		if err = os.MkdirAll(jc.defaultDir, 0700); err != nil {
			return err
		}
		path = filepath.Join(jc.defaultDir, configFileName)
	}

	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	log.WithFields(logrus.Fields{"file": path}).Debug("Writing config file")
	return ioutil.WriteFile(path, b, 0700)
}

func (jc *jsonConfig) Get(section, key string) (string, bool, error) {
	c, err := jc.load()
	if err != nil {
		return "", false, err
	}

	v, ok := c[section][key]
	if !ok {
		return "", false, nil
	}

	return string(v), true, nil
}

func (jc *jsonConfig) Set(section, key, value string) error {
	c, err := jc.load()
	if err != nil {
		return err
	}

	if _, ok := c[section]; !ok {
		c[section] = map[string]boolOrString{}
	}

	c[section][key] = boolOrString(value)
	if err = jc.save(c); err != nil {
		return err
	}

	return nil
}

func (jc *jsonConfig) Sections() ([]string, error) {
	c, err := jc.load()
	if err != nil {
		return nil, err
	}

	sections := []string{}
	for k := range c {
		sections = append(sections, k)
	}

	return sections, nil
}

func (jc *jsonConfig) Section(section string) (map[string]string, error) {
	c, err := jc.load()
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for key, val := range c[section] {
		result[key] = string(val)
	}

	return result, nil
}
