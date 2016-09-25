package config

import (
	"os"
	"path/filepath"

	xdg "github.com/casimir/xdg-go"
)

const (
	configFileName = "config.json"
)

var (
	app = xdg.App{Name: "cardigann"}
)

func GetConfigPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if f, exists := fileExists(cwd, configFileName); exists {
		return f, nil
	}

	if configDir := os.Getenv("CONFIG_DIR"); configDir != "" {
		return filepath.Join(configDir, configFileName), nil
	}

	return app.ConfigPath(configFileName), nil
}

func GetDefinitionDirs() ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dirs := []string{
		filepath.Join(cwd, "definitions"),
		app.ConfigPath("definitions"),
	}

	return append(dirs, app.SystemConfigPaths("definitions")...), nil
}

func fileExists(f ...string) (string, bool) {
	full := filepath.Join(f...)
	if _, err := os.Stat(full); os.IsNotExist(err) {
		return full, false
	}
	return full, true
}
