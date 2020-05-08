package systemd

// Reload triggers a reload of the systemd daemon configuration.
func Reload() error {
	return systemctl("daemon-reload")
}
