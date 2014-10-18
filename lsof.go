package procspy

// `lsof` executing implementation

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	lsofFields = "cPn" // parseLSOF() depends on the order
)

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
func parseLSOF(out string) ([]ConnectionProc, error) {
	res := []ConnectionProc{}
	cp := ConnectionProc{}
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
			cp.LocalAddress = net.ParseIP(localAddr)
			p, err := strconv.Atoi(localPort)
			if err != nil {
				return nil, err
			}
			cp.LocalPort = uint16(p)

			remoteAddr, remotePort, err := net.SplitHostPort(addresses[1])
			if err != nil {
				return nil, fmt.Errorf("invalid remote address field: %v", err)
			}
			cp.RemoteAddress = net.ParseIP(remoteAddr)
			p, err = strconv.Atoi(remotePort)
			if err != nil {
				return nil, err
			}
			cp.RemotePort = uint16(p)

			if cp.PID != 0 {
				res = append(res, cp)
			}
		default:
			return nil, fmt.Errorf("unexpected lsof field: %v in %#v", field, value)
		}
	}
	return res, nil
}
