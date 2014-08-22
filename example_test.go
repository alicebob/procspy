package procspy_test

import (
	"fmt"

	"github.com/alicebob/procspy"
)

func Example() {
	procs, err := procspy.Spy()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connections:\n")
	for _, p := range procs {
		fmt.Printf(" - %v\n", p)
	}
}
