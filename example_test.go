package procspy_test

import (
	"fmt"

	"github.com/alicebob/procspy"
)

func Example() {
	lookupProcesses := true
	cs, err := procspy.Connections(lookupProcesses)
	if err != nil {
		panic(err)
	}

	fmt.Printf("TCP Connections:\n")
	for _, c := range cs {
		fmt.Printf(" - %v\n", c)
	}
}
