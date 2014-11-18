package procspy

// lsof-executing implementation.

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	lsofFields = "cPi" // parseLSOF() depends on the order
)

// parseLsof parses lsof out with `-F cPn` argument.
//
// Format description: the first letter is the type of record, records are
// newline seperated, the record starting with 'p' (pid) is a new processid.
// There can be multiple connections for the same 'p' record in which case the
// 'p' is not repeated.
//
// For example, this is one process with two listens and one connection:
//
//   p13100
//   cmpd
//   PTCP
//   n[::1]:6600
//   PTCP
//   n127.0.0.1:6600
//   PTCP
//   n[::1]:6600->[::1]:50992
//
func parseLSOF(out string) (Procs, error) {
	var (
		res = Procs{}
		cp  = Proc{}
	)
	for _, line := range strings.Split(out, "\n") {
		if len(line) <= 1 {
			continue
		}

		var (
			field = line[0]
			value = line[1:]
		)
		switch field {
		case 'p':
			pid, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid 'p' field in lsof output: %#v", value)
			}
			cp.PID = uint(pid)

		case 'i':
			inode, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid 'i' field in lsof output: %#v", value)
			}
			res[uint64(inode)] = cp

		case 'c':
			cp.Name = value

		default:
			return nil, fmt.Errorf("unexpected lsof field: %v in %#v", field, value)
		}
	}

	return res, nil
}
