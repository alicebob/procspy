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
