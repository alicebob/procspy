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

// Connection is a (TCP) connection.
type Connection struct {
	Transport     string
	LocalAddress  net.IP
	LocalPort     uint16
	RemoteAddress net.IP
	RemotePort    uint16
	Inode         uint64
}

// Proc is a single connection with PID/process name.
type Proc struct {
	PID  uint
	Name string
}

// Procs is a mapping from Inode to Proc
type Procs map[uint64]Proc

// ConnIter is returned by Connections().
type ConnIter interface {
	Next() *Connection
}

// Connections returns all established (TCP) connections.
// No need to be root to run this.
func Connections() (ConnIter, error) {
	return cbConnections()
}

// Processes makes the inode -> process map. You can combine these with the
// output from Connections().
func Processes() (Procs, error) {
	return cbProcesses()
}
