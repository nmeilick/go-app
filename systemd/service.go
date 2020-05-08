package systemd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Service contains a systemd service configuration.
type Service struct {
	Name        string
	Description string
	After       []string
	Before      []string

	Type       string
	PrivateTmp bool

	ExecStart        string
	ExecReload       string
	ExecStop         string
	WorkingDirectory string

	User  string
	Group string

	Restart                  string
	RestartSec               string
	RestartPreventExitStatus []string

	// Text, if non-empty, contains the service file content to be written.
	Text string
}

// New service returns an empty service.
func NewService(name string) *Service {
	return &Service{
		Name: name,
	}
}

// NewSimpleService returns a service definition with defaults for a
// simple-type service.
func NewSimpleService(name string) *Service {
	return &Service{
		Name:        name,
		Description: "",
		After:       []string{"network.target"},
		Before:      nil,

		Type:       "simple",
		PrivateTmp: true,

		ExecStart:        "",
		ExecReload:       "/bin/kill -s HUP $MAINPID",
		ExecStop:         "",
		WorkingDirectory: "",

		User:  "",
		Group: "",

		Restart:                  "",
		RestartSec:               "1s",
		RestartPreventExitStatus: []string{"SIGTERM", "SIGINT"},
	}
}

// Validate returns an error if the service contains errors.
func (svc *Service) Validate() error {
	switch {
	case svc == nil:
		return errors.New("is nil")
	case svc.Name == "":
		return errors.New("name is empty")
	case svc.Type == "":
		return errors.New("type is empty")
	}
	return nil
}

// File returns the default location of the service file.
func (svc *Service) File() string {
	return filepath.Join(defaultServicePath, "system", svc.Name+".service")
}

// Files returns all existing service files of the service.
func (svc *Service) Files() (files []string) {
	if svc.Validate() != nil {
		return
	}

	for _, p := range SearchPaths() {
		file := filepath.Join(p, "system", svc.Name+".service")
		if stat, err := os.Lstat(file); err == nil && !stat.IsDir() {
			files = append(files, file)
		}
	}
	return
}

// Exists returns true if a service file exists.
func (svc *Service) Exists() bool {
	for _, f := range svc.Files() {
		if _, err := os.Stat(f); err == nil {
			return true
		}
	}
	return false
}

// Create writes a new service file and triggers a daemon reload.
func (svc *Service) Create() error {
	if err := svc.Validate(); err != nil {
		return errors.Wrap(err, "invalid service")
	}

	file := svc.File()

	tmp, err := ioutil.TempFile(filepath.Dir(file), filepath.Base(file))
	if err != nil {
		return err
	}
	_ = tmp.Chmod(0644)

	defer func() {
		if _, err := os.Stat(tmp.Name()); err == nil || !os.IsNotExist(err) {
			tmp.Close()
			os.Remove(tmp.Name())
		}
	}()

	var text string
	if svc.Text == "" {
		lines := []string{"[Unit]"}

		if svc.Description != "" {
			lines = append(lines, "Description="+svc.Description)
		}
		if len(svc.Before) > 0 {
			lines = append(lines, "Before="+strings.Join(svc.Before, " "))
		}
		if len(svc.After) > 0 {
			lines = append(lines, "After="+strings.Join(svc.After, " "))
		}
		lines = append(lines, "")
		lines = append(lines, "[Service]")
		lines = append(lines, "Type="+svc.Type)
		if svc.PrivateTmp {
			lines = append(lines, "PrivateTmp=true")
		}

		if svc.User != "" {
			lines = append(lines, "User="+svc.User)
		} else {
			lines = append(lines, "User=root")
		}

		if svc.Group != "" {
			lines = append(lines, "Group="+svc.Group)
		} else {
			lines = append(lines, "Group=root")
		}

		if svc.ExecStart != "" {
			lines = append(lines, "ExecStart="+svc.ExecStart)
		}
		if svc.ExecReload != "" {
			lines = append(lines, "ExecReload="+svc.ExecReload)
		}
		if svc.ExecStop != "" {
			lines = append(lines, "ExecStop="+svc.ExecStop)
		}
		if svc.WorkingDirectory != "" {
			lines = append(lines, "WorkingDirectory="+svc.WorkingDirectory)
		} else if svc.ExecStart != "" {
			lines = append(lines, "WorkingDirectory="+filepath.Dir(svc.ExecStart))
		}

		if svc.Restart != "" {
			lines = append(lines, "Restart="+svc.Restart)
			if svc.RestartSec != "" {
				lines = append(lines, "RestartSec="+svc.RestartSec)
			}
			if len(svc.RestartPreventExitStatus) != 0 {
				lines = append(lines, "RestartPreventExitStatus="+strings.Join(svc.RestartPreventExitStatus, " "))
			}
		}

		lines = append(lines, "")
		lines = append(lines, "[Install]")
		lines = append(lines, "WantedBy=multi-user.target")

		text = strings.Join(lines, "\n")
	} else {
		text = svc.Text
	}

	if _, err := fmt.Fprintln(tmp, text); err != nil {
		return err
	}

	if err := os.Rename(tmp.Name(), file); err != nil {
		return err
	}

	if err := Reload(); err != nil {
		return err
	}
	return systemctl("enable", svc.Name)
}

// Delete removes all existing service files and triggers a daemon reload.
func (svc *Service) Delete() error {
	var errtext []string
	for _, f := range svc.Files() {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			errtext = append(errtext, err.Error())
		}
	}

	if err := Reload(); err != nil {
		errtext = append(errtext, err.Error())
	}

	if len(errtext) > 0 {
		return errors.New(strings.Join(errtext, "; "))
	}

	return nil
}

// Start triggers a start of the service.
func (svc *Service) Start() error {
	return systemctl("start", svc.Name)
}

// Stop triggers a stop of the service.
func (svc *Service) Stop() error {
	return systemctl("stop", svc.Name)
}

// Reload triggers a reload of the service.
func (svc *Service) Reload() error {
	return systemctl("reload", svc.Name)
}
