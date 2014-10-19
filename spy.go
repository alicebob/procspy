// Package procspy lists TCP connections, and optionally tries to find the
// owning processes.
// Works on Linux (via /proc) and Darwin (via `lsof -i` and `netstat`).
// You'll need root to use Processes().
package procspy

import (
	"net"
)

// Connection is a (TCP) connection.
type Connection struct {
	Transport     string
	LocalAddress  net.IP
	LocalPort     uint16
	RemoteAddress net.IP
	RemotePort    uint16
	inode         uint64
}

// ConnectionProc is a single connection with PID/process name.
type ConnectionProc struct {
	Connection
	PID  uint
	Name string
}
