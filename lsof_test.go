package procspy

import (
	"net"
	"reflect"
	"testing"
)

func TestLSOFParsing(t *testing.T) {
	// List of lsof -> expected entries
	for in, expected := range map[string][]ConnectionProc{
		// Single connection
		"p25196\n" +
			"ccello-app\n" +
			"PTCP\n" +
			"n127.0.0.1:48094->127.0.0.1:4039\n" +
			"PTCP\n" +
			"n*:4040\n": []ConnectionProc{
			{
				Connection: Connection{
					Transport:     "tcp",
					LocalAddress:  net.ParseIP("127.0.0.1"),
					LocalPort:     48094,
					RemoteAddress: net.ParseIP("127.0.0.1"),
					RemotePort:    4039,
				},
				PID:  25196,
				Name: "cello-app",
			},
		},

		// Only listen()s.
		"cdhclient\n" +
			"PUDP\n" +
			"n*:68\n" +
			"PUDP\n" +
			"n*:38282\n" +
			"PUDP\n" +
			"n*:40625\n": []ConnectionProc{},

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
			"n192.168.2.111:33963->192.168.2.71:22\n": []ConnectionProc{
			{
				Connection: Connection{
					Transport:     "tcp",
					LocalAddress:  net.ParseIP("::1"),
					LocalPort:     6600,
					RemoteAddress: net.ParseIP("::1"),
					RemotePort:    50992,
				},
				PID:  13100,
				Name: "mpd",
			},
			{
				Connection: Connection{
					Transport:     "tcp",
					LocalAddress:  net.ParseIP("2003:45:2b57:8900:1869:2947:f942:aba7"),
					LocalPort:     55711,
					RemoteAddress: net.ParseIP("2a00:1450:4008:c01::11"),
					RemotePort:    443,
				},
				PID:  14612,
				Name: "chromium",
			},
			{
				Connection: Connection{
					Transport:     "tcp",
					LocalAddress:  net.ParseIP("192.168.2.111"),
					LocalPort:     37158,
					RemoteAddress: net.ParseIP("192.0.72.2"),
					RemotePort:    80,
				},
				PID:  14612,
				Name: "chromium",
			},
			{
				Connection: Connection{
					Transport:     "tcp",
					LocalAddress:  net.ParseIP("192.168.2.111"),
					LocalPort:     44013,
					RemoteAddress: net.ParseIP("54.229.241.196"),
					RemotePort:    80,
				},
				PID:  14612,
				Name: "chromium",
			},
			{
				Connection: Connection{
					Transport:     "tcp",
					LocalAddress:  net.ParseIP("192.168.2.111"),
					LocalPort:     56385,
					RemoteAddress: net.ParseIP("74.201.105.31"),
					RemotePort:    443,
				},
				PID:  14612,
				Name: "chromium",
			},
			{
				Connection: Connection{
					Transport:     "tcp",
					LocalAddress:  net.ParseIP("192.168.2.111"),
					LocalPort:     33963,
					RemoteAddress: net.ParseIP("192.168.2.71"),
					RemotePort:    22,
				},
				PID:  21356,
				Name: "ssh",
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
