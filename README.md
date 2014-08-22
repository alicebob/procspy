Go module to list all TCP connections with PID and processname.

Uses /proc directly where /proc is available, with a fallback to `lsof -i` when
not.

Works for IPv4 and IPv6 TCP connections. Only active connections are listed, ports where something is only listening are skipped. Connections where the owning process is unknown are also skipped.

If you want all connections you'll need to run this as root.

Status:
-------

Tested on Linux and Darwin (10.9). Should work everywhere where `lsof` is available.

Install:
--------

`go install`

Usage:
------

`list, err := procspy.Spy()`

(See ./example\_test.go)

If you want you can call `procspy.SpyProc()` or `procspy.SpyLSOF()` directly.

``` go

package main

import (
	"fmt"

	"github.com/alicebob/procspy"
)

func main() {
	procs, err := procspy.Spy()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connections:\n")
	for _, p := range procs {
		fmt.Printf(" - %v\n", p)
	}
}
```
