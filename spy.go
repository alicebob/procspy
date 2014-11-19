// Package procspy lists TCP connections, and optionally tries to find the
// owning processes. Works on Linux (via /proc) and Darwin (via `lsof -i` and
// `netstat`). You'll need root to use Processes().
package procspy

import (
	"net"
)

const (
	tcpEstablished = 1 // according to /include/net/tcp_states.h
)

// Connection is a (TCP) connection. The Proc struct might not be filled in.
type Connection struct {
	Transport     string
	LocalAddress  net.IP
	LocalPort     uint16
	RemoteAddress net.IP
	RemotePort    uint16
	inode         uint64
	Proc
}

// Proc is a single connection with PID/process name.
type Proc struct {
	PID  uint
	Name string
}

// ConnIter is returned by Connections().
type ConnIter interface {
	Next() *Connection
}

// Connections returns all established (TCP) connections.
// No need to be root to run this. If processes is true we'll try to lookup the
// process owning the connection. You will need to run this as root to find all
// processes.
func Connections(processes bool) (ConnIter, error) {
	return cbConnections(processes)
}
