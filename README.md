Go module to list all TCP connections with PID and processname.

Works by reading /proc on Linux, and by executing `netstat` and `lsof` on Darwin.
not.

Works for IPv4 and IPv6 TCP connections. Only active connections are listed, ports where something is only listening are skipped. Connections where the owning process is unknown are also skipped.

If you want all processes you'll need to run this as root.

Status:
-------

Tested on Linux and Darwin (10.9).

Install:
--------

`go install`

Usage:
------

`conns, err := procspy.Connections()`

`connProcs, err := procspy.Processes()`

(See ./example\_test.go)

``` go

package main

import (
	"fmt"

	"github.com/alicebob/procspy"
)

func main() {
	cs, err := procspy.Connections()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connections:\n")
	for _, c := range cs {
		fmt.Printf(" - %v\n", c)
	}
}
```
