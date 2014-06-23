package procspy


import (
	"net"
)

type ConnProc struct {
	Protocol  string
	LocalAddr net.IP
	LocalPort uint16
	PID       uint
	Name      string
}

/*
func main() {
	// HASNCACHE NcacheReload = 1
	// C.print_lsof()
	for {
		fmt.Printf("Go Go:\n")
		for _, p := range lsof() {
			fmt.Printf(" - %v\n", p)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
*/
