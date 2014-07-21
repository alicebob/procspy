package procspy

// /proc based implementation

import (
	"bufio"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const procRoot = "/proc"

// SpyProc uses /proc directly to make the connection list.
func SpyProc() ([]ConnProc, error) {
	// A map of inode -> pid
	inodes, err := walkProcPid()
	if err != nil {
		return nil, err
	}

	res := []ConnProc{}
	for _, procFile := range []string{
		procRoot + "/net/tcp",
		procRoot + "/net/tcp6",
	} {
		fh, err := os.Open(procFile)
		if err != nil {
			// File might not be there if IPv{4,6} is not supported.
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
					// Remote address is zero. This is a listen entry.
					continue
				}
				res = append(res, ConnProc{
					Transport:  "tcp",
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
	dirNames, err := fh.Readdirnames(-1)
	fh.Close()
	if err != nil {
		return nil, err
	}
	procmap := map[uint64]uint{}
	var stat syscall.Stat_t
	for _, dirName := range dirNames {
		pid, err := strconv.ParseUint(dirName, 10, 0)
		if err != nil {
			// Not a number, so not a PID subdir.
			continue
		}

		fdBase := procRoot + "/" + dirName + "/fd/"
		dfh, err := os.Open(fdBase)
		if err != nil {
			// Process is be gone by now, or we don't have access.
			continue
		}
		fdNames, err := dfh.Readdirnames(-1)
		dfh.Close()
		if err != nil {
			continue
		}
		for _, fdName := range fdNames {
			// We want sockets only
			// Direct use of Stat() to save garbage.
			err = syscall.Stat(fdBase+fdName, &stat)
			if err != nil {
				continue
			}
			if stat.Mode&syscall.S_IFMT != syscall.S_IFSOCK {
				continue
			}
			procmap[stat.Ino] = uint(pid)
		}
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
func parseTransport(r io.Reader) []transport {
	res := []transport{}
	scanner := bufio.NewScanner(r)
	for i := 0; scanner.Scan(); i++ {
		if i == 0 {
			// Skip header
			continue
		}
		// Fields are:
		//  'sl local_address rem_address st tx_queue rx_queue tr tm->when retrnsmt uid timeout inode <more>'
		fields := strings.Fields(scanner.Text())
		if len(fields) < 10 {
			continue
		}

		localAddress, localPort, err := scanAddress(fields[1])
		if err != nil {
			continue
		}

		remoteAddress, remotePort, err := scanAddress(fields[2])
		if err != nil {
			continue
		}

		uid, err := strconv.Atoi(fields[7])
		if err != nil {
			continue
		}

		inode, err := strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
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

// scanAddress parses 'A12CF62E:E4D7' to the address and port.
// Handles IPv4 and IPv6 addresses.
// The address part are big endian 32 bit ints, hex encoded. Since net.IP is a
// byte slice we just decode the hex and flip the bytes in every group of 4.
func scanAddress(in string) (net.IP, uint16, error) {
	parts := strings.Split(in, ":")
	if len(parts) != 2 {
		return nil, 0, errors.New("invalid address:port")
	}

	// Network address is big endian. Can be either ipv4 or ipv6.
	address, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, 0, err
	}
	// reverse every 4 byte-sequence.
	for i := 0; i < len(address); i += 4 {
		address[i], address[i+3] = address[i+3], address[i]
		address[i+1], address[i+2] = address[i+2], address[i+1]
	}

	// Port number
	port, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return nil, 0, err
	}

	return net.IP(address), uint16(port), err
}

// procName does a pid->name lookup
func procName(pid uint) (string, error) {
	fh, err := os.Open(procRoot + "/" + strconv.FormatUint(uint64(pid), 10) + "/comm")
	if err != nil {
		return "", err
	}
	name := make([]byte, 1024)
	l, err := fh.Read(name)
	fh.Close()
	if err != nil {
		return "", err
	}
	if l < 2 {
		return "", nil
	}
	// drop trailing "\n"
	return string(name[:l-1]), nil
}
