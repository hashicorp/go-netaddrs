package netaddrs

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
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
	retIPAddrs, err := IPAddrs("google.com", log.New(ioutil.Discard, "netaddrs: ", 0))
	if err != nil {
		t.Fatalf("Error resolving DNS name to IP addresses. %s", err)
	}
	err = validIPAddrs(retIPAddrs)
	if err != nil {
		t.Fatalf("IP address invalid. %s", err)
	}
}

func TestIPAddrsDNSFail(t *testing.T) {
	_, err := IPAddrs("invalidDNSname", log.New(ioutil.Discard, "netaddrs: ", 0))
	if err == nil {
		t.Fatalf("Expected error on invalid DNS name")
	}
}

func TestIPAddrsCustomExecutable(t *testing.T) {
	type testCase struct {
		name      string
		cmd       string
		expectErr string
	}

	run := func(t *testing.T, tc testCase) {
		retIPAddrs, err := IPAddrs(tc.cmd, log.New(ioutil.Discard, "netaddrs: ", 0))
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
		if err = validIPAddrs(retIPAddrs); err != nil {
			t.Fatalf("IP address invalid. %s", err)
		}
	}

	testcases := []testCase{
		{
			name: "custom executable without args",
			cmd:  "exec=sample_scripts/ipaddrs_valid_without_args.sh",
		},
		{
			name: "custom executable with args and same line output",
			cmd:  "exec=sample_scripts/ipaddrs_valid_with_args.sh same-line",
		},
		{
			name: "custom executable with args and multi line ipv6 addresses",
			cmd:  "exec=sample_scripts/ipaddrs_valid_with_args.sh multi-line",
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
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
