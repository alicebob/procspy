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
	fmt.Printf("TCP Connections:\n")
	for _, c := range cs {
		fmt.Printf(" - %v\n", c)
	}

	ps, err := procspy.Processes(cs)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Processes:\n")
	for _, p := range ps {
		fmt.Printf(" - %v\n", p)
	}
}
