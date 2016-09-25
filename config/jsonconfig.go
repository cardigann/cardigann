package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type jsonConfig struct {
	path string
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

// NewJSONConfig creates a new json config backed on the given file
func NewJSONConfig(f string) (Config, error) {
	return &jsonConfig{f}, nil
}

func (jc *jsonConfig) load() (jsonConfigMap, error) {
	config := jsonConfigMap{}

	data, err := ioutil.ReadFile(jc.path)
	if os.IsNotExist(err) {
		return config, nil
	} else if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func (jc *jsonConfig) save(c jsonConfigMap) error {
	if err := os.MkdirAll(filepath.Dir(jc.path), 0700); err != nil {
		return err
	}

	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(jc.path, b, 0700)
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
