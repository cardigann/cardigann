package config

import (
	"encoding/json"

	"github.com/shibukawa/configdir"
)

const (
	configFileName = "config.json"
)

type jsonConfig struct {
	configDirs configdir.ConfigDir
}

func (jc *jsonConfig) load() (configMap, error) {
	config := configMap{}
	folder := jc.configDirs.QueryFolderContainsFile(configFileName)
	if folder != nil {
		// log.Printf("Loading config from %s", folder.Path)
		data, err := folder.ReadFile(configFileName)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, err
		}
	}
	return config, nil
}

func (jc *jsonConfig) save(c configMap) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	var folder *configdir.Config
	if folder = jc.configDirs.QueryFolderContainsFile(configFileName); folder == nil {
		folders := jc.configDirs.QueryFolders(configdir.Global)
		folder = folders[0]
	}

	// log.Printf("Saving config to %s/%s", folder.Path, configFileName)
	return folder.WriteFile(configFileName, b)
}

func (jc *jsonConfig) Get(section, key string) (string, bool, error) {
	c, err := jc.load()
	if err != nil {
		return "", false, err
	}

	// log.Printf("config.Get[%s,%s] => %#v", section, key, c)
	v, ok := c[section][key]
	return v, ok, nil
}

func (jc *jsonConfig) Set(section, key, value string) error {
	c, err := jc.load()
	if err != nil {
		return err
	}

	c[section][key] = value

	err = jc.save(c)
	if err != nil {
		return err
	}

	// log.Printf("config.Set[%s,%s,%s] => %#v", section, key, value, c)
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

	return c[section], nil
}
