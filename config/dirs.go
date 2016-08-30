package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/shibukawa/configdir"
)

func configDir() (*configdir.ConfigDir, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cd := configdir.New("cardigann", "cardigann")
	cd.LocalPath = cwd

	return &cd, nil
}

func Dirs() ([]string, error) {
	cd, err := configDir()
	if err != nil {
		return nil, err
	}

	configDirs := []string{}
	for _, folder := range cd.QueryFolders(configdir.All) {
		configDirs = append(configDirs, folder.Path)
	}

	return configDirs, nil
}

func DefaultDir() (string, error) {
	cd, err := configDir()
	if err != nil {
		return "", err
	}

	folder := cd.QueryFolders(configdir.Local)
	return folder[0].Path, nil
}

func Find(filename string, dirs []string) (string, error) {
	for i := len(dirs) - 1; i >= 0; i-- {
		path := filepath.Join(dirs[i], filename)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", errors.New("File doesn't exist")
}
