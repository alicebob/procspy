package procspy

import (
	"reflect"
	"testing"
)

func TestLSOFParsing(t *testing.T) {
	// List of lsof -> expected entries
	for in, expected := range map[string][]ConnProc{
		// Single connection
		"p25196\n" +
			"ccello-app\n" +
			"PTCP\n" +
			"n127.0.0.1:48094->127.0.0.1:4039\n" +
			"PTCP\n" +
			"n*:4040\n": []ConnProc{
			{
				Transport:  "tcp",
				LocalAddr:  "127.0.0.1",
				LocalPort:  uint16(48094),
				RemoteAddr: "127.0.0.1",
				RemotePort: uint16(4039),
				PID:        25196,
				Name:       "cello-app",
			},
		},

		// Only listen()s.
		"cdhclient\n" +
			"PUDP\n" +
			"n*:68\n" +
			"PUDP\n" +
			"n*:38282\n" +
			"PUDP\n" +
			"n*:40625\n": []ConnProc{},

		// A bunch
		"p13100\n" +
			"cmpd\n" +
			"PTCP\n" +
			"n[::1]:6600\n" +
			"PTCP\n" +
			"n127.0.0.1:6600\n" +
			"PTCP\n" +
			"n[::1]:6600->[::1]:50992\n" +
			"p14612\n" +
			"cchromium\n" +
			"PTCP\n" +
			"n[2003:45:2b57:8900:1869:2947:f942:aba7]:55711->[2a00:1450:4008:c01::11]:443\n" +
			"PTCP\n" +
			"n192.168.2.111:37158->192.0.72.2:80\n" +
			"PTCP\n" +
			"n192.168.2.111:44013->54.229.241.196:80\n" +
			"PTCP\n" +
			"n192.168.2.111:56385->74.201.105.31:443\n" +
			"p21356\n" +
			"cssh\n" +
			"PTCP\n" +
			"n192.168.2.111:33963->192.168.2.71:22\n": []ConnProc{
			{
				Transport:  "tcp",
				LocalAddr:  "::1",
				LocalPort:  uint16(6600),
				RemoteAddr: "::1",
				RemotePort: uint16(50992),
				PID:        13100,
				Name:       "mpd",
			},
			{
				Transport:  "tcp",
				LocalAddr:  "2003:45:2b57:8900:1869:2947:f942:aba7",
				LocalPort:  uint16(55711),
				RemoteAddr: "2a00:1450:4008:c01::11",
				RemotePort: uint16(443),
				PID:        14612,
				Name:       "chromium",
			},
			{
				Transport:  "tcp",
				LocalAddr:  "192.168.2.111",
				LocalPort:  uint16(37158),
				RemoteAddr: "192.0.72.2",
				RemotePort: uint16(80),
				PID:        14612,
				Name:       "chromium",
			},
			{
				Transport:  "tcp",
				LocalAddr:  "192.168.2.111",
				LocalPort:  uint16(44013),
				RemoteAddr: "54.229.241.196",
				RemotePort: uint16(80),
				PID:        14612,
				Name:       "chromium",
			},
			{
				Transport:  "tcp",
				LocalAddr:  "192.168.2.111",
				LocalPort:  uint16(56385),
				RemoteAddr: "74.201.105.31",
				RemotePort: uint16(443),
				PID:        14612,
				Name:       "chromium",
			},
			{
				Transport:  "tcp",
				LocalAddr:  "192.168.2.111",
				LocalPort:  uint16(33963),
				RemoteAddr: "192.168.2.71",
				RemotePort: uint16(22),
				PID:        21356,
				Name:       "ssh",
			},
		},
	} {
		got, err := parseLSOF(in)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !reflect.DeepEqual(expected, got) {
			t.Errorf("Expected:\n %#v\nGot:\n %#v\n", expected, got)
		}
	}
}
