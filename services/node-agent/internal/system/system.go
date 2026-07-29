package system

import (
	"os"
	"runtime"
)

// Info contains static information about the host system.
type Info struct {
	Hostname     string `json:"hostname"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	GoVersion    string `json:"go_version"`
}

// GetInfo returns static information about the host.
func GetInfo() Info {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return Info{
		Hostname:     hostname,
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		GoVersion:    runtime.Version(),
	}
}