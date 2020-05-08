package systemd

// Path where service files will be created.
const defaultServicePath = "/etc/systemd"

// Paths where systemd is looking for service files.
var searchPaths = []string{
	defaultServicePath,
	"/usr/lib/systemd",
	"/lib/systemd",
}

// SearchPaths returns a list of paths where systemd is looking for service files.
func SearchPaths() []string {
	return searchPaths
}
