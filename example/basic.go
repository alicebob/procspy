package main

import (
	"fmt"

	"github.com/alicebob/procspy"
)

func main() {
	fmt.Printf("Go Go:\n")
	for _, p := range procspy.Spy() {
		fmt.Printf(" - %v\n", p)
	}
}
