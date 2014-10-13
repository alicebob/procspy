package procspy

import (
	"net"
	"reflect"
	"strings"
	"testing"
)

func TestTransport4(t *testing.T) {
	testString := strings.NewReader(`  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode                                                     
   0: 00000000:A6C0 00000000:0000 0A 00000000:00000000 00:00000000 00000000   105        0 5107 1 ffff8800a6aaf040 100 0 0 10 0                      
   1: 00000000:006F 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 5084 1 ffff8800a6aaf740 100 0 0 10 0                      
   2: 0100007F:0019 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 10550 1 ffff8800a729b780 100 0 0 10 0                     
   3: A12CF62E:E4D7 57FC1EC0:01BB 01 00000000:00000000 02:000006FA 00000000  1000        0 639474 2 ffff88007e75a740 48 4 26 10 -1                   
`)
	res := ParseTransport(testString)
	expected := []transport{
		transport{
			LocalAddress: net.IP{0x0, 0x0, 0x0, 0x0}, LocalPort: 0xa6c0,
			RemoteAddress: net.IP{0x0, 0x0, 0x0, 0x0}, RemotePort: 0x0,
			UID:   105,
			Inode: 5107,
		},
		transport{
			LocalAddress: net.IP{0x0, 0x0, 0x0, 0x0}, LocalPort: 0x006f,
			RemoteAddress: net.IP{0x0, 0x0, 0x0, 0x0}, RemotePort: 0x0,
			UID:   0,
			Inode: 5084,
		},
		transport{
			LocalAddress: net.IP{0x7f, 0x0, 0x0, 0x01}, LocalPort: 0x0019,
			RemoteAddress: net.IP{0x0, 0x0, 0x0, 0x0}, RemotePort: 0x0,
			UID:   0,
			Inode: 10550,
		},
		transport{
			LocalAddress: net.IP{0x2e, 0xf6, 0x2c, 0xa1}, LocalPort: 0xe4d7,
			RemoteAddress: net.IP{0xc0, 0x1e, 0xfc, 0x57}, RemotePort: 0x01bb,
			UID:   1000,
			Inode: 639474,
		},
	}

	if len(res) != 4 {
		t.Errorf("Wanted 4")
	}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("transport 4 error. Got\n%+v\nExpected\n%+v\n", res, expected)
	}

}

func TestTransport6(t *testing.T) {
	// Abridged copy of my /proc/net/tcp6
	testString := strings.NewReader(` sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000000000000000000000000000:19C8 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 23661201 1 ffff880103fb4800 100 0 0 10 -1
   8: 4500032000BE692B8AE31EBD919D9D10:D61C 5014002A080805400000000015100000:01BB 01 00000000:00000000 02:00000045 00000000  1000        0 36856710 2 ffff88010b796080 22 4 30 8 7
`)

	res := ParseTransport(testString)
	expected := []transport{
		transport{
			LocalAddress: net.IP(make([]byte, 16)), LocalPort: 0x19c8,
			RemoteAddress: net.IP(make([]byte, 16)), RemotePort: 0x0,
			UID:   0,
			Inode: 23661201,
		},
		transport{
			LocalAddress: net.IP{
				0x20, 0x03, 0, 0x45,
				0x2b, 0x69, 0xbe, 0x00,
				0xbd, 0x1e, 0xe3, 0x8a,
				0x10, 0x9d, 0x9d, 0x91,
			},
			LocalPort: 0xd61c,
			RemoteAddress: net.IP{
				0x2a, 0x00, 0x14, 0x50,
				0x40, 0x05, 0x08, 0x08,
				0, 0, 0, 0,
				0, 0, 0x10, 0x15,
			},
			RemotePort: 0x01bb,
			UID:        1000,
			Inode:      36856710,
		},
	}

	if len(res) != 2 {
		t.Errorf("Wanted 2")
	}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("transport 6 error. Got\n%+v\nExpected\n%+v\n", res, expected)
	}

}
