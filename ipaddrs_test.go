package netaddrs

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"strings"
	"testing"
)

func validIPAddrs(ipaddrs []net.IPAddr) error {
	for _, addr := range ipaddrs {
		ipaddr := net.ParseIP(addr.IP.String())
		if ipaddr == nil {
			return fmt.Errorf("invalid IP address format: %s.", addr)
		}
	}
	return nil
}

func TestIPAddrsDNS(t *testing.T) {
	retIPAddrs, err := IPAddrs(context.Background(), "google.com", noopLogger{})
	if err != nil {
		t.Fatalf("Error resolving DNS name to IP addresses. %s", err)
	}
	err = validIPAddrs(retIPAddrs)
	if err != nil {
		t.Fatalf("IP address invalid. %s", err)
	}
}

type noopLogger struct{}

func (l noopLogger) Debug(msg string, args ...interface{}) {
}

func TestIPAddrsDNSFail(t *testing.T) {
	_, err := IPAddrs(context.Background(), "invalidDNSname", noopLogger{})
	if err == nil {
		t.Fatalf("Expected error on invalid DNS name")
	}
}

func TestIPAddrsCustomExecutable(t *testing.T) {
	type testCase struct {
		name        string
		cmd         string
		expectErr   string
		expectedIPs []net.IPAddr
	}

	run := func(t *testing.T, tc testCase) {
		actual, err := IPAddrs(context.Background(), tc.cmd, noopLogger{})
		if tc.expectErr != "" {
			if err == nil {
				t.Fatalf("Expected error return, got nil")
			}
			if actual := err.Error(); !strings.Contains(actual, tc.expectErr) {
				t.Fatalf("Expected error %q to contain %q", actual, tc.expectErr)
			}
			return
		}
		if err != nil {
			t.Fatalf("Error retrieving IP addrs on running executable. %s", err)
		}
		if err = validIPAddrs(actual); err != nil {
			t.Fatalf("IP address invalid. %s", err)
		}
		if !reflect.DeepEqual(actual, tc.expectedIPs) {
			t.Fatalf("expected %v, got %v", tc.expectedIPs, actual)
		}
	}

	testcases := []testCase{
		{
			name: "custom executable without args",
			cmd:  "exec=sample_scripts/ipaddrs_valid_without_args.sh",
			expectedIPs: []net.IPAddr{
				ipAddr("172.25.41.79"),
				ipAddr("172.25.16.77"),
				ipAddr("172.25.42.80"),
			},
		},
		{
			name: "custom executable with args and same line output",
			cmd:  "exec=sample_scripts/ipaddrs_valid_with_args.sh same-line",
			expectedIPs: []net.IPAddr{
				ipAddr("172.25.41.79"),
				ipAddr("172.25.16.77"),
				ipAddr("172.25.42.80"),
			},
		},
		{
			name: "custom executable with args and multi line ipv6 addresses",
			cmd:  "exec=sample_scripts/ipaddrs_valid_with_args.sh multi-line",
			expectedIPs: []net.IPAddr{
				ipAddr("2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
				{IP: net.ParseIP("fe80::1ff:fe23:4567:890a"), Zone: "3"},
				ipAddr("2001:db8::"),
			},
		},
		{
			name:      "custom executable with invalid ip address output",
			cmd:       "exec=sample_scripts/ipaddrs_invalid1.sh",
			expectErr: "executable returned invalid IP address: 172.25.16.77:8080",
		},
		{
			name:      "custom executable returned error",
			cmd:       "exec=sample_scripts/ipaddrs_invalid2.sh",
			expectErr: "executable failed with exit code 1: ERROR! No Consul servers found.",
		},
		{
			name:        "custom executable with stderr",
			cmd:         "exec=sample_scripts/ipaddrs_valid_with_stderr.sh",
			expectedIPs: []net.IPAddr{ipAddr("172.10.1.123")},
		},
		{
			name:      "custom executable with no output",
			cmd:       "exec=sample_scripts/no_output.sh",
			expectErr: "executable returned no output",
		},
		{
			name:      "custom executable not found",
			cmd:       "exec=sample_scripts/not_found.sh",
			expectErr: "not_found.sh: no such file or directory",
		},
		{
			name:      "custom executable not executable",
			cmd:       "exec=./ipaddrs_test.go",
			expectErr: "permission denied",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}

func ipAddr(v string) net.IPAddr {
	return net.IPAddr{IP: net.ParseIP(v)}
}
