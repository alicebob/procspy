package procspy
import (
	"fmt"
	"testing"
	"net"
	"strings"
	"reflect"
)

func TestTransport4(t *testing.T) {
	testString := strings.NewReader(`sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode                                                     
   0: 00000000:A6C0 00000000:0000 0A 00000000:00000000 00:00000000 00000000   105        0 5107 1 ffff8800a6aaf040 100 0 0 10 0                      
   1: 00000000:006F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 5084 1 ffff8800a6aaf740 100 0 0 10 0                      
   2: 0100007F:0019 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 10550 1 ffff8800a729b780 100 0 0 10 0                     
   3: A12CF62E:E4D7 57FC1EC0:01BB 01 00000000:00000000 02:000006FA 00000000  1000        0 639474 2 ffff88007e75a740 48 4 26 10 -1                   
`)
	res := parseTransport4(testString)
	expected := []transport{
		transport{
			localAddress:net.IP{0x0, 0x0, 0x0, 0x0}, localPort:0xa6c0,
			remoteAddress:net.IP{0x0, 0x0, 0x0, 0x0}, remotePort:0x0,
			uid:105,
			inode:5107,
		},
		transport{
			localAddress:net.IP{0x0, 0x0, 0x0, 0x0}, localPort:0x006f,
			remoteAddress:net.IP{0x0, 0x0, 0x0, 0x0}, remotePort:0x0,
			uid:0,
			inode:5084,
		},
		transport{
			localAddress:net.IP{0x7f, 0x0, 0x0, 0x01}, localPort:0x0019,
			remoteAddress:net.IP{0x0, 0x0, 0x0, 0x0}, remotePort:0x0,
			uid:0,
			inode:10550,
		},
		transport{
			localAddress:net.IP{0x2e, 0xf6, 0x2c, 0xa1}, localPort:0xe4d7,
			remoteAddress:net.IP{0xc0, 0x1e, 0xfc, 0x57}, remotePort:0x01bb,
			uid:1000,
			inode:639474,
		},
	}

	if len(res) != 4 {
		t.Errorf("Wanted 4")
	}
	if ! reflect.DeepEqual(res, expected) {
		t.Errorf("transport 4 error. Got\n%+v\nExpected\n%+v\n", res, expected)
	}

}

/*
func TestWalkproc(t *testing.T) {
	walkprocpid()
}
*/

func TestSpy(t *testing.T) {
	for _, t := range Spy() {
		fmt.Printf("A proc: %+v\n", t)
	}
}
