package app

import (
	"errors"
	"os"
	"path/filepath"
)

// ConfigFiles returns a slice of possible config file locations.
func (p *Properties) ConfigFiles() (files []string) {
	if len(p.Name) > 0 {
		for _, path := range p.SearchPaths() {
			files = append(files, filepath.Join(path, p.Name+".conf"))
		}
	}
	return
}

// FindConfig searches the application search paths for an existing config file
// and returns the first one, otherwise an error is returned.
func (p *Properties) FindConfig() (string, error) {
	if p.Name == "" {
		return "", errors.New("Application name not set")
	}

	for _, f := range p.ConfigFiles() {
		if stat, err := os.Stat(f); err == nil && !stat.IsDir() {
			return f, nil
		}
	}

	return "", errors.New("no config found")
}
