package config

import (
	"os"
	"strings"
)

type Config interface {
	Get(section, key string) (string, bool, error)
	Set(section, key, value string) error
	Sections() ([]string, error)
	Section(section string) (map[string]string, error)
}

// IsSectionEnabled returns true if a section has an enabled=true
func IsSectionEnabled(section string, c Config) bool {
	v, _, err := c.Get(section, "enabled")
	if err != nil {
		return false
	}

	return v == "ok" || v == "true"
}

// GetDefault gets a value from the config, with a default if it doesn't exist
func GetDefault(section, key, defaultVal string, c Config) (string, error) {
	val, ok, err := c.Get(section, key)
	if err != nil || !ok {
		return defaultVal, err
	}
	return val, nil
}

const (
	GlobalConfigSection = "global"
)

// GetGlobalConfig gets a value from the Environment, the global config or a default fallback value if
// neither exists. If the key is "test", "CARDIGANN_TEST" is checked in the environment
func GetGlobalConfig(key, defaultVal string, c Config) (string, error) {
	envVal := os.Getenv("CARDIGANN_" + strings.ToUpper(key))
	if envVal != "" {
		return envVal, nil
	}
	val, ok, err := c.Get(GlobalConfigSection, key)
	if err != nil || !ok {
		return defaultVal, err
	}
	return val, nil
}
