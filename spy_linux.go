package procspy

import (
	"bytes"
	"sync"
)

// sync.Pool turns out cheaper than keeping a freelist.
var bufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 5000))
	},
}

var cbConnections = func() ([]Connection, error) {
	var c []Connection
	buf := bufPool.Get().(*bytes.Buffer)
	for _, procFile := range []string{
		procRoot + "/net/tcp",
		procRoot + "/net/tcp6",
	} {
		buf.Reset()
		if err := readFile(procFile, buf); err != nil {
			// File might not be there if IPv{4,6} is not supported.
			continue
		}
		// Only read established connections.
		c = append(c, parseTransport(buf.String(), tcpEstablished)...)
	}
	bufPool.Put(buf)
	return c, nil
}

var cbProcesses = func(c []Connection) ([]ConnectionProc, error) {
	return procProcesses(c), nil
}
