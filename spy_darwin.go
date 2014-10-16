package procspy

import (
	"os/exec"
)

const (
	netstatBinary = "netstat"
	lsofBinary    = "lsof"
)

// Connections returns all (TCP) connections.
// No need to be root to run this.
func Connections() ([]Connection, error) {
	out, err := exec.Command(
		netstatBinary,
		"-n", // no number resolving
		"-W", // Wide output
		// "-l", // full IPv6 addresses // What does this do?
		"-p", "tcp", // only TCP
	).CombinedOutput()
	if err != nil {
		// log.Printf("lsof error: %s", err)
		return nil, err
	}
	return parseDarwinNetstat(string(out)), nil
}

// Processes returns the list of connections with processes.
// You need to be root to run this.
func Processes() ([]ConnectionProc, error) {
	out, err := exec.Command(
		lsofBinary,
		"-i",       // only Internet files
		"-n", "-P", // no number resolving
		"-w",             // no warnings
		"-F", lsofFields, // \n based output of only the fields we want.
	).CombinedOutput()
	if err != nil {
		// log.Printf("lsof error: %s", err)
		return nil, err
	}
	return parseLSOF(string(out))
}
