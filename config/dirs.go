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

func GetDefinitionDirs() []string {
	dirs := []string{}

	if cwd, err := os.Getwd(); err == nil {
		dirs = append(dirs, filepath.Join(cwd, "definitions"))
	}

	dirs = append(dirs, app.ConfigPath("definitions"))

	if configDir := os.Getenv("CONFIG_DIR"); configDir != "" {
		dirs = append(dirs, filepath.Join(configDir, "definitions"))
	}

	return append(dirs, app.SystemConfigPaths("definitions")...)
}

func GetCachePath(file string) string {
	return app.CachePath(file)
}

func fileExists(f ...string) (string, bool) {
	full := filepath.Join(f...)
	if _, err := os.Stat(full); os.IsNotExist(err) {
		return full, false
	}
	return full, true
}
