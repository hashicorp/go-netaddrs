package netaddrs

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"testing"
)

func validIPAddrs(ipaddrs []net.IPAddr) error {
	for _, addr := range ipaddrs {
		ipaddr := net.ParseIP(addr.String())
		if ipaddr == nil {
			return fmt.Errorf("Invalid IP address format: %s.", addr)
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
	testcases := []struct {
		name      string
		cmd       string
		expectErr bool
	}{
		{"custom executable without args", "exec=sample_scripts/ipaddrs_valid_without_args.sh", false},
		{"custom executable with args and same line output", "exec=sample_scripts/ipaddrs_valid_with_args.sh same-line", false},
		{"custom executable with args and multi line output", "exec=sample_scripts/ipaddrs_valid_with_args.sh multi-line", false},
		{"custom executable with invalid ip address output", "exec=sample_scripts/ipaddrs_invalid1", true},
		{"custom executable returned error", "exec=sample_scripts/ipaddrs_invalid2.sh", true},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			retIPAddrs, err := IPAddrs(tc.cmd, log.New(ioutil.Discard, "netaddrs: ", 0))
			if tc.expectErr {
				if err == nil {
					t.Fatalf("Expected error on running executable.")
				}
			} else {
				if err != nil {
					t.Fatalf("Error retrieving IP addrs on running executable. %s", err)
				}
				err = validIPAddrs(retIPAddrs)
				if err != nil {
					t.Fatalf("IP address invalid. %s", err)
				}

			}
		})
	}
}
