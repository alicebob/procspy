package procspy

import (
	"os/exec"
)

const (
	netstatBinary = "netstat"
	lsofBinary    = "lsof"
)

// Connections returns all established (TCP) connections.
// No need to be root to run this.
var cbConnections = func() (ConnIter, error) {
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

	f := fixtConnIter(parseDarwinNetstat(string(out)))
	return &f, nil
}

// Processes returns map with Inode->Process. You can combine this with the
// output from Connections.
// You need to be root to find all processes.
var cbProcesses = func() (Procs, error) {
	out, err := exec.Command(
		lsofBinary,
		"-i",       // only Internet files
		"-n", "-P", // no number resolving
		"-w",             // no warnings
		"-F", lsofFields, // \n based output of only the fields we want.
	).CombinedOutput()
	if err != nil {
		return nil, err
	}
	return parseLSOF(string(out))
}
