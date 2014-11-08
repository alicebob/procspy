package procspy

// /proc based implementation

import (
	"bytes"
	"net"
	"os"
	"strconv"
	"syscall"
)

const (
	procRoot = "/proc"
)

// procProcesses gives all processes for the given connections. It is used by
// the linux version of Processes().
func procProcesses(conn []Connection) []ConnectionProc {
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
			// conn only has 'Established' connections.
			// if net.IP(tp.RemoteAddress).IsUnspecified() {
			// // Remote address is zero. This is a listen entry.
			// continue
			// }
			res = append(res, ConnectionProc{
				// Connection: tp,
				Connection: Connection{
					Transport:     tp.Transport,
					LocalAddress:  tp.LocalAddress,
					LocalPort:     tp.LocalPort,
					RemoteAddress: tp.RemoteAddress,
					RemotePort:    tp.RemotePort,
					inode:         tp.inode,
				},
				PID:  pid,
				Name: name,
			})
		}
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

// parseTransport parses /proc/net/{tcp,udp}{,6} files.
// It will filter out all rows not in wantedState.
func parseTransport(s []byte, wantedState uint) []Connection {
	// The file format is well-known, so we use some specialized versions of
	// std lib functions to speed things up a bit.

	res := make([]Connection, 0, len(s)/149/10) // heuristic. Lines are about 150 chars long, and say we have 10% established.

	// Lines are:
	// '  sl local_address rem_address st tx_queue rx_queue tr tm->when retrnsmt uid timeout inode <more>'
	// '  0: 00000000:0FC9 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 11276449 1 ffff8801029607c0 100 0 0 10 0'

	// Skip header
	s = nextLine(s)

	var (
		local, remote, state, inode []byte
	)
	for {
		if len(s) == 0 {
			break
		}
		_, s = nextField(s) // 'sl' column
		local, s = nextField(s)
		remote, s = nextField(s)
		state, s = nextField(s)
		if parseHex(state) != wantedState {
			s = nextLine(s)
			continue
		}
		_, s = nextField(s) // 'tx_queue' column
		_, s = nextField(s) // 'rx_queue' column
		_, s = nextField(s) // 'tr' column
		_, s = nextField(s) // 'uid' column
		_, s = nextField(s) // 'timeout' column
		inode, s = nextField(s)

		t := Connection{}
		t.LocalAddress, t.LocalPort = scanAddress(local)
		t.RemoteAddress, t.RemotePort = scanAddress(remote)
		t.inode = parseDec(inode)
		res = append(res, t)

		s = nextLine(s)

	}
	return res
}

// scanAddress parses 'A12CF62E:00AA' to the address/port. Handles IPv4 and
// IPv6 addresses.  The address is a big endian 32 bit ints, hex encoded. We
// just decode the hex and flip the bytes in every group of 4.
func scanAddress(in []byte) (net.IP, uint16) {
	col := bytes.IndexByte(in, ':')
	if col == -1 {
		return nil, 0
	}
	// Network address is big endian. Can be either ipv4 or ipv6.
	address := hexDecode(in[:col])
	// reverse every 4 byte-sequence.
	for i := 0; i < len(address); i += 4 {
		address[i], address[i+3] = address[i+3], address[i]
		address[i+1], address[i+2] = address[i+2], address[i+1]
	}
	port := parseHex(in[col+1:])
	return net.IP(address), uint16(port)
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

func nextField(s []byte) ([]byte, []byte) {
	// Skip whitespace.
	for i, b := range s {
		if b != ' ' {
			s = s[i:]
			break
		}
	}
	// Up until the next non-space field.
	for i, b := range s {
		if b == ' ' {
			return s[:i], s[i:]
		}
	}
	return nil, nil
}

func nextLine(s []byte) []byte {
	for i, b := range s {
		if b == '\n' {
			return s[i+1:]
		}
	}
	return nil
}

// readFile reads a /proc file info a buffer.
func readFile(filename string, buf *bytes.Buffer) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	_, err = buf.ReadFrom(f)
	f.Close()
	return err
}

// Simplified copy of strconv.ParseUint(16).
func parseHex(s []byte) uint {
	n := uint(0)
	for i := 0; i < len(s); i++ {
		n *= 16
		n += uint(fromHexChar(s[i]))
	}
	return n
}

// Simplified copy of strconv.ParseUint(10).
func parseDec(s []byte) uint64 {
	n := uint64(0)
	for _, c := range s {
		n *= 10
		n += uint64(c - '0')
	}
	return n
}

// hexDecode and fromHexChar are taken from encoding/hex.
func hexDecode(src []byte) []byte {
	if len(src)%2 == 1 {
		return nil
	}

	dst := make([]byte, len(src)/2)
	for i := 0; i < len(src)/2; i++ {
		a := fromHexChar(src[i*2])
		b := fromHexChar(src[i*2+1])
		dst[i] = (a << 4) | b
	}

	return dst
}

// fromHexChar converts a hex character into its value.
func fromHexChar(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}

	return 0
}
