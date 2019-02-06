package app

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var argv0 = func() string { return os.Args[0] }()

// Properties stores info about the running application.
type Properties struct {
	Executable string
	Dir        string
	Name       string
}

// isFile returns true if the given path points to an existing file.
func isFile(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil && !stat.Mode().IsDir())
}

// GetExecutable tries to find the full path of the command the program was
// started with or an empty string if the path could not be determined.
func GetExecutable() string {
	if argv0 != "" {
		if strings.Contains(argv0, "/") || strings.Contains(argv0, "\\") {
			if dir, err := os.Getwd(); err == nil {
				file := filepath.Clean(filepath.Join(dir, argv0))
				if isFile(file) {
					return file
				}
			}
		} else if file, err := exec.LookPath(argv0); err == nil {
			if file, err = filepath.Abs(file); err == nil {
				file = filepath.Clean(file)
				if isFile(file) {
					return file
				}
			}
		}
	}

	if file, err := os.Executable(); err == nil {
		if file, err = filepath.Abs(file); err == nil {
			return file
		}
	}

	if argv0 != "" {
		if dir, err := os.Getwd(); err == nil {
			return filepath.Clean(filepath.Join(dir, argv0))
		}
	}

	return ""
}

// GetProperties returns info about the running application.
//
func GetProperties() (*Properties, error) {
	p := Properties{
		Executable: GetExecutable(),
	}

	if p.Executable == "" {
		return &p, errors.New("Could not find my executable")
	}
	p.Dir = filepath.Dir(p.Executable)

	if argv0 != "" {
		// Get the base path and remove any extensions
		p.Name = strings.Split(filepath.Base(argv0), ".")[0]
	}

	if p.Name == "" {
		p.Name = strings.Split(filepath.Base(p.Executable), ".")[0]
		if p.Name == "" {
			p.Name = "invalid"
			return &p, errors.New("Could not determine application name")
		}
	}
	return &p, nil
}
