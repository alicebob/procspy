package procspy

import (
	"bufio"
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

func Spy() []ConnProc {
	// A map of inode -> pid
	inodes, err := walkProcPid()
	if err != nil {
		fmt.Printf("walkProcPid err: %v\n", err)
		return nil
	}

	res := []ConnProc{}

	fh, err := os.Open("/proc/net/tcp")
	defer fh.Close()
	if err != nil {
		fmt.Printf("Open err: %v\n", err)
	}
	if err == nil {
		// fmt.Printf("Go parseTransport4\n")
		for _, tp := range parseTransport4(fh) {
			// fmt.Printf("Tp: %+v\n", tp)
			if pid, ok := inodes[tp.inode]; ok {
				res = append(res, ConnProc{
					Protocol:  "tcp",
					LocalAddr: tp.localAddress,
					LocalPort: tp.localPort,
					PID:       pid,
					Name:      "unknown yet", // <-- todo
				})
			}
		}
	}
	return res
}

func walkProcPid() (map[uint64]uint, error) {
	// Walk over all /proc entries (numerical ones, those are PIDs), and see if their ./fd/* files link to 'socket[...]' 'files'.
	// Returns a map from socket id to pid
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
			// Not a number, so not a pid subdir.
			continue
		}
		// fmt.Printf("Pid: %v\n", pid)

		dfh, err := os.Open(procRoot + "/" + dirName + "/fd")
		if err != nil {
			// process can be gone by now, or we don't have access.
			// fmt.Printf("Skip %v: %v\n", pid, err)
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
			// fmt.Printf(" Proc: %v: %v\n", pid, procFd.Name())

			// We want sockets only
			stat, err := os.Stat(procRoot + "/" + dirName + "/fd/" + procFd.Name())
			if err != nil {
				continue
			}
			if stat.Mode()&os.ModeSocket == 0 {
				continue
			}
			// fmt.Printf(" Stat: %v : %v\n", pid, stat.Name())
			sys, ok := stat.Sys().(*syscall.Stat_t)
			if !ok {
				panic("Weird result from stat.Sys()")
			}
			// fmt.Printf(" Inode: %v : %v\n", pid, sys.Ino)
			procmap[sys.Ino] = pid

			/*
				linkName, err := os.Readlink(procRoot + "/" + pid + "/fd/" + procFd.Name())
				if err != nil {
					fmt.Printf("Readlink err: %v\n", err)
					continue
				}
				fmt.Printf(" Link: %v: %v\n", pid, linkName)
			*/
			/*
				eval, err := filepath.EvalSymlinks(procRoot + "/" + name + "/fd/" + procFd.Name())
				if err != nil {
					fmt.Printf("EvalSymlinks err: %v\n", err)
					continue
				}
				fmt.Printf(" Link: %v: %v\n", name, eval)
			*/
			// fmt.Printf(" Proc: %v: %v\n", name, procFd.Name())
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

// parseTransport4 parses /proc/net/{tcp,udp} files
var fieldRe = regexp.MustCompile(`\s+`)

func parseTransport4(r io.Reader) []transport {
	res := []transport{}
	scanner := bufio.NewScanner(r)
	for i := 0; scanner.Scan(); i++ {
		if i == 0 {
			// fmt.Printf("Header: %s\n", scanner.Text())
			continue
		}
		// fmt.Printf("Line: %s\n", scanner.Text())
		// Fields are 'sl local_address rem_address st tx_queue rx_queue tr tm->when retrnsmt uid timeout inode <more>'
		line := strings.TrimSpace(scanner.Text())
		fields := fieldRe.Split(line, -1)
		if len(fields) < 10 {
			// log.Error("Invalid format")
			fmt.Printf("Invalid format")
			continue
		}

		localAddress, localPort, err := scanAddress4(fields[1])
		if err != nil {
			fmt.Printf("err: %s\n", err)
			continue
		}

		remoteAddress, remotePort, err := scanAddress4(fields[2])
		if err != nil {
			fmt.Printf("err: %s\n", err)
			continue
		}

		var uid int
		_, err = fmt.Sscanf(fields[7], "%d", &uid)
		if err != nil {
			fmt.Printf("err: %s\n", err)
			continue
		}

		var inode uint64
		_, err = fmt.Sscanf(fields[9], "%d", &inode)
		if err != nil {
			fmt.Printf("err: %s\n", err)
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

// scanAddress4 parses things like 'A12CF62E:E4D7' to their address/port
func scanAddress4(in string) (net.IP, uint16, error) {
	parts := strings.Split(in, ":")
	if len(parts) != 2 {
		return nil, 0, errors.New("invalid addres:port")
	}

	// network (big) endian
	address := make([]byte, 4)
	_, err := fmt.Sscanf(parts[0], "%2X%2X%2X%2X", &address[3], &address[2], &address[1], &address[0])
	if err != nil {
		return nil, 0, err
	}

	var port uint16
	_, err = fmt.Sscanf(parts[1], "%X", &port)
	if err != nil {
		return nil, 0, err
	}

	return net.IP(address), port, err
}
