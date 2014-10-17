package procspy

import (
	"strconv"
)

const (
	tcpEstablished = 1 // according to /include/net/tcp_states.h
)

// Connections returns all established (TCP) connections.
// No need to be root to run this.
func Connections() ([]Connection, error) {
	var c []Connection
	for _, pc := range procConnections() {

		if pc.State != tcpEstablished {
			continue
		}
		// // Skip zero addresses. They are listening sockets.
		// if pc.LocalAddress.IsUnspecified() || pc.RemoteAddress.IsUnspecified() {
		// continue
		// }

		c = append(c, Connection{
			Transport:     "tcp",
			LocalAddress:  pc.LocalAddress.String(),
			LocalPort:     strconv.Itoa(int(pc.LocalPort)),
			RemoteAddress: pc.RemoteAddress.String(),
			RemotePort:    strconv.Itoa(int(pc.RemotePort)),
		})
	}
	return c, nil
}

// Processes returns the list of Connections, and tries to find the process
// which handles the connection.
// Only connections for which we found a process are returned.
// You need to be root to find all processes.
func Processes() ([]ConnectionProc, error) {
	return procProcesses(procConnections()), nil
}
