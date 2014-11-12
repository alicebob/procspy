package procspy

import (
	"bytes"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 5000))
	},
}

type pnConnIter struct {
	pn  *ProcNet
	buf *bytes.Buffer
}

func (c *pnConnIter) Next() *Connection {
	n := c.pn.Next()
	if n == nil {
		// Done!
		bufPool.Put(c.buf)
	}
	return n
}

// cbConnections sets Connections()
var cbConnections = func() (ConnIter, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	readFile(procRoot+"/net/tcp", buf)
	readFile(procRoot+"/net/tcp6", buf)
	return &pnConnIter{
		pn:  NewProcNet(buf.Bytes(), tcpEstablished),
		buf: buf,
	}, nil
}

// cbProcesses sets Processes()
var cbProcesses = func() (Procs, error) {
	return walkProcPid()
}
