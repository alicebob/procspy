Go module to list all TCP connections by port->process.

Uses /proc directly where available, with a fallback to the `lsof` binary.

Works for IPv4 and IPv6
TCP connections. Only active connections are listed, ports where something is
listening are skipped.

Status:
-------

Works on Linux and Darwin (10.9, Mavericks)

Install:
--------

`go install`

Usage:
------

`procspy.Spy()`

(See ./example/)

If you want you can directly call `procspy.SpyProc()` or `procspy.SpyLSOF()`

``` go

package main

import (
	"fmt"

	"github.com/alicebob/procspy"
)

func main() {
	fmt.Printf("Go Go:\n")
	procs, err := procspy.Spy()
	if err != nil {
		panic(err)
	}
	for _, p := range procs {
		fmt.Printf(" - %v\n", p)
	}
}
```
