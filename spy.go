package procspy

import (
	"os"
)

// ConnProc is a single connection with PID/process name.
type ConnProc struct {
	Protocol   string
	LocalAddr  string
	LocalPort  uint16
	RemoteAddr string
	RemotePort uint16
	PID        uint
	Name       string
}

// Spy returns the current []ConnProc list.
// It will use /proc if that's available, otherwise it'll fallback to `lsof -i`
func Spy() ([]ConnProc, error) {
	if _, err := os.Stat("/proc"); err == nil {
		return SpyProc()
	}
	return SpyLSOF()
}
