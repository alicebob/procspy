package procspy

// /proc based implementation

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
	"syscall"
)

const procRoot = "/proc"

// SpyProc uses /proc directly to make the connection list.
func SpyProc() ([]ConnProc, error) {
	// A map of inode -> pid
	inodes, err := walkProcPid()
	if err != nil {
		// fmt.Printf("walkProcPid err: %v\n", err)
		return nil, err
	}

	res := []ConnProc{}
	for _, procFile := range []string{
		"/proc/net/tcp",
		"/proc/net/tcp6",
	} {
		fh, err := os.Open(procFile)
		if err != nil {
			// fmt.Printf("Open err: %v\n", err)
			// Might not be there is IPv{4,6} is not supported.
			continue
		}
		defer fh.Close()
		for _, tp := range parseTransport(fh) {
			if pid, ok := inodes[tp.inode]; ok {
				name, err := procName(pid)
				if err != nil {
					// Process might be gone by now
					continue
				}
				if tp.remoteAddress.IsUnspecified() {
					// remote address is zero. This is a listen entry.
					continue
				}
				res = append(res, ConnProc{
					Protocol:   "tcp",
					LocalAddr:  tp.localAddress.String(),
					LocalPort:  tp.localPort,
					RemoteAddr: tp.remoteAddress.String(),
					RemotePort: tp.remotePort,
					PID:        pid,
					Name:       name,
				})
			}
		}
	}
	return res, nil
}

func walkProcPid() (map[uint64]uint, error) {
	// Walk over all /proc entries (numerical ones, those are PIDs), and see if
	// their ./fd/* files are symlink to sockets.
	// Returns a map from socket id ('inode`) to PID.
	// Will return an error if /proc/ isn't there.
	fh, err := os.Open(procRoot)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	dirNames, err := fh.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	procmap := map[uint64]uint{}
	for _, dirName := range dirNames {
		var pid uint
		_, err = fmt.Sscanf(dirName, "%d", &pid)
		if err != nil {
			// Not a number, so not a PID subdir.
			continue
		}

		dfh, err := os.Open(procRoot + "/" + dirName + "/fd")
		if err != nil {
			// process is be gone by now, or we don't have access.
			continue
		}
		fds, err := dfh.Readdir(-1)
		if err != nil {
			dfh.Close()
			continue
		}
		for _, procFd := range fds {
			if procFd.Mode()&os.ModeSymlink == 0 {
				continue
			}

			// We want sockets only
			stat, err := os.Stat(procRoot + "/" + dirName + "/fd/" + procFd.Name())
			if err != nil {
				continue
			}
			if stat.Mode()&os.ModeSocket == 0 {
				continue
			}
			sys, ok := stat.Sys().(*syscall.Stat_t)
			if !ok {
				panic("Weird result from stat.Sys()")
			}
			procmap[sys.Ino] = pid
		}
		dfh.Close()
	}
	return procmap, nil
}

// transport are found in /proc/net/{tcp,udp}{,6} files
type transport struct {
	localAddress  net.IP
	localPort     uint16
	remoteAddress net.IP
	remotePort    uint16
	uid           int
	inode         uint64
}

// parseTransport parses /proc/net/{tcp,udp}{,6} files
var fieldRe = regexp.MustCompile(`\s+`)

func parseTransport(r io.Reader) []transport {
	res := []transport{}
	scanner := bufio.NewScanner(r)
	for i := 0; scanner.Scan(); i++ {
		if i == 0 {
			continue
		}
		// Fields are:
		//  'sl local_address rem_address st tx_queue rx_queue tr tm->when retrnsmt uid timeout inode <more>'
		line := strings.TrimSpace(scanner.Text())
		fields := fieldRe.Split(line, -1)
		if len(fields) < 10 {
			// fmt.Printf("Invalid format")
			continue
		}

		localAddress, localPort, err := scanAddress(fields[1])
		if err != nil {
			// fmt.Printf("err: %s\n", err)
			continue
		}

		remoteAddress, remotePort, err := scanAddress(fields[2])
		if err != nil {
			// fmt.Printf("err: %s\n", err)
			continue
		}

		var uid int
		_, err = fmt.Sscanf(fields[7], "%d", &uid)
		if err != nil {
			// fmt.Printf("err: %s\n", err)
			continue
		}

		var inode uint64
		_, err = fmt.Sscanf(fields[9], "%d", &inode)
		if err != nil {
			// fmt.Printf("err: %s\n", err)
			continue
		}
		t := transport{
			localAddress:  localAddress,
			localPort:     localPort,
			remoteAddress: remoteAddress,
			remotePort:    remotePort,
			uid:           uid,
			inode:         inode,
		}
		res = append(res, t)

	}
	return res
}

// scanAddress parses things like 'A12CF62E:E4D7' to their address/port.
// Deals with IPv4 and IPv6 addresses.
// The address part are big endian 32 bit ints, hex encoded. Since net.IP is a
// byte slice we just decode the hex and flip the bytes in every group of 4.
func scanAddress(in string) (net.IP, uint16, error) {
	parts := strings.Split(in, ":")
	if len(parts) != 2 {
		return nil, 0, errors.New("invalid address:port")
	}

	// network (big) endian. Can be either ipv4 or ipv6
	address, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, 0, err
	}
	// reverse every 4 byte-sequence.
	for i := 0; i < len(address); i += 4 {
		address[i], address[i+3] = address[i+3], address[i]
		address[i+1], address[i+2] = address[i+2], address[i+1]
	}

	var port uint16
	_, err = fmt.Sscanf(parts[1], "%X", &port)
	if err != nil {
		return nil, 0, err
	}

	return net.IP(address), port, err
}

// procName does a pid->name lookup
func procName(pid uint) (string, error) {
	fh, err := os.Open(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return "err1", err
	}
	name := make([]byte, 1024)
	l, err := fh.Read(name)
	if err != nil {
		return "err", err
	}
	if l < 2 {
		return "", nil
	}
	// drop trailing "\n"
	return string(name[:l-1]), nil
}
