Go bindings for `lsof -i` (connection -> local process name)

Status:
-------

Works on Linux and Darwin (10.9, Mavericks)

Install:
--------

`make` (will use sudo to install things in /usr/local/lib/procspy)

Usage:
------

```
	import (
		"github.com/alicebob/procspy"
		"fmt"
	)

	func main() {
		for _, p := range procspy.Spy() {
			fmt.Printf("- %v\n", p)
		}
	}
```
