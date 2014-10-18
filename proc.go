package procspy

// /proc based implementation

import (
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const (
	procRoot = "/proc"
)

// procProcesses gives all processes for the given connections. It is used by
// the linux version of Processes().
func procProcesses(conn []transport) []ConnectionProc {
	// A map of inode -> pid
	inodes, err := walkProcPid()
	if err != nil {
		return nil
	}

	res := []ConnectionProc{}
	for _, tp := range conn {
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
			res = append(res, ConnectionProc{
				Connection: Connection{
					Transport:     "tcp",
					LocalAddress:  tp.localAddress.String(),
					LocalPort:     strconv.Itoa(int(tp.localPort)),
					RemoteAddress: tp.remoteAddress.String(),
					RemotePort:    strconv.Itoa(int(tp.remotePort)),
				},
				PID:  pid,
				Name: name,
			})
		}
	}
	return res
}

// connections gives all TCP IPv{4,6} connections as found in
// /proc/net/tcp{,6}.  It is used by the linux version of Processes().
func procConnections() []transport {
	var res []transport
	for _, procFile := range []string{
		procRoot + "/net/tcp",
		procRoot + "/net/tcp6",
	} {
		c, err := ioutil.ReadFile(procFile)
		if err != nil {
			// File might not be there if IPv{4,6} is not supported.
			continue
		}
		res = append(res, parseTransport(string(c))...)
	}
	return res
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
			// Direct use of syscall.Stat() to save garbage.
			err = syscall.Stat(fdBase+fdName, &stat)
			if err != nil {
				continue
			}
			// We want sockets only
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
	state         int
	localAddress  net.IP
	localPort     uint16
	remoteAddress net.IP
	remotePort    uint16
	uid           int
	inode         uint64
}

// parseTransport parses /proc/net/{tcp,udp}{,6} files
func parseTransport(s string) []transport {
	res := make([]transport, 0, 10)
	for i, line := range strings.Split(s, "\n") {
		if i == 0 {
			// Skip header
			continue
		}
		// Fields are:
		// 'sl local_address rem_address st tx_queue rx_queue tr tm->when retrnsmt uid timeout inode <more>'
		fields := procNetFields(line)
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

		state, err := strconv.ParseInt(fields[3], 16, 32)
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
		res = append(res, transport{
			state:         int(state),
			localAddress:  localAddress,
			localPort:     localPort,
			remoteAddress: remoteAddress,
			remotePort:    remotePort,
			uid:           uid,
			inode:         inode,
		})

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

// Copy of the standard stings.FieldsFunc(), but just for our tcp lines.
func procNetFields(s string) []string {
	// We know there are 18 fields.
	n := 24 // buffer if the file changes.

	a := make([]string, n)
	na := 0
	fieldStart := -1 // Set to -1 when looking for start of field.
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			if fieldStart >= 0 {
				a[na] = s[fieldStart:i]
				na++
				fieldStart = -1
			}
		} else if fieldStart == -1 {
			fieldStart = i
		}
	}
	if fieldStart >= 0 { // Last field might end at EOF.
		a[na] = s[fieldStart:]
	}
	return a
}
