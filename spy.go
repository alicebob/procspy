package procspy

import (
	"os"
)

type ConnProc struct {
	Protocol   string
	LocalAddr  string // net.IP
	LocalPort  uint16
	RemoteAddr string // net.IP
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
