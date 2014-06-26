package procspy

// `lsof` executing implementation

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	lsofBinary   = "lsof"
	lsofFields   = "cPn" // parseLSOF() depends on the order
	lsofInterval = 250 * time.Millisecond
)

// SpyLSOF executes `lsof` to get the connection list.
func SpyLSOF() ([]ConnProc, error) {
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

// parseLsof parses lsof out with `-F cPn` argument.
// Format description: the first letter is the type of record, records are
// newline seperated, the record starting with 'p' (pid) is a new processid.
// There can be multiple connections for the same 'p' record in which case the
// 'p' is not repeated.  For example, this is one process with two listens and
// one connection:
//   p13100
//   cmpd
//   PTCP
//   n[::1]:6600
//   PTCP
//   n127.0.0.1:6600
//   PTCP
//   n[::1]:6600->[::1]:50992
func parseLSOF(out string) ([]ConnProc, error) {
	res := []ConnProc{}
	cp := ConnProc{}
	for _, line := range strings.Split(out, "\n") {
		if len(line) <= 1 {
			continue
		}
		field := line[0]
		value := line[1:]
		switch field {
		case 'p':
			pid, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid 'p' field in lsof output: %#v", value)
			}
			cp.PID = uint(pid)
		case 'c':
			cp.Name = value
		case 'P':
			cp.Transport = strings.ToLower(value)
		case 'u':
			// uid
		case 'n':
			// 'n' is the last field, with '-F cPn'
			// format examples:
			// "192.168.2.111:44013->54.229.241.196:80"
			// "[2003:45:2b57:8900:1869:2947:f942:aba7]:55711->[2a00:1450:4008:c01::11]:443"
			// "*:111" <- a listen
			addresses := strings.SplitN(value, "->", 2)
			if len(addresses) != 2 {
				// That's a listen entry.
				continue
			}
			localAddr, localPort, err := net.SplitHostPort(addresses[0])
			if err != nil {
				return nil, fmt.Errorf("invalid local address field: %v", err)
			}
			cp.LocalAddr = localAddr
			_, err = fmt.Sscanf(localPort, "%d", &cp.LocalPort)
			if err != nil {
				return nil, fmt.Errorf("invalid local port number: %v", localPort)
			}

			remoteAddr, remotePort, err := net.SplitHostPort(addresses[1])
			if err != nil {
				return nil, fmt.Errorf("invalid remote address field: %v", err)
			}
			cp.RemoteAddr = remoteAddr
			_, err = fmt.Sscanf(remotePort, "%d", &cp.RemotePort)
			if err != nil {
				return nil, fmt.Errorf("invalid remote port number: %v", remotePort)
			}

			if cp.PID != 0 {
				res = append(res, cp)
			}
		default:
			return nil, fmt.Errorf("unexpected lsof field: %v in %#v", field, value)
		}
	}
	return res, nil
}
