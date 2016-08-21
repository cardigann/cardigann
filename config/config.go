package config

import (
	"os"

	"github.com/shibukawa/configdir"
)

type configMap map[string]map[string]string

type Config interface {
	Get(section, key string) (string, bool, error)
	Set(section, key, value string) error
	Sections() ([]string, error)
	Section(section string) (map[string]string, error)
}

func NewConfig() (Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cd := configdir.New("cardigann", "cardigann")
	cd.LocalPath = cwd

	return &jsonConfig{cd}, nil
}

func IsEnabled(section string, c Config) bool {
	v, ok, err := c.Get(section, "enabled")
	if err != nil {
		return false
	}

	return v == "ok" || v == "true" || !ok
}
