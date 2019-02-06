package app

import (
	"errors"
	"os"
	"path/filepath"
)

// FindConfig searches the application search paths for an existing config file
// and returns the first one, otherwise an error is returned.
func (p *Properties) FindConfig() (string, error) {
	if p.Name == "" {
		return "", errors.New("Application name not set")
	}

	for _, path := range p.SearchPaths() {
		file := filepath.Join(path, p.Name+".conf")
		if stat, err := os.Stat(file); err == nil && !stat.IsDir() {
			return file, nil
		}
	}

	return "", errors.New("no config found")
}
