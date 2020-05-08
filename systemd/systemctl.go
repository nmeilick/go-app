package systemd

import "github.com/nmeilick/go-run"

// systemctl executes the external systemctl command with the given arguments.
func systemctl(args ...string) error {
	return run.Run("systemctl", args...).Error()
}
