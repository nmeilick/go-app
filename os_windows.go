package app

import (
	"os"
	"path/filepath"
)

// SearchPaths returns application specific search paths, from most specific
// to least specific.
func (p *Properties) SearchPaths() (paths []string) {
	if p.Name != "" {
		if path := os.Getenv("USERPROFILE"); path != "" {
			paths = append(paths, filepath.Join(path, "."+p.Name))
		}
	}

	if p.Dir != "" {
		paths = append(paths, p.Dir)
	}
	return
}
