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
