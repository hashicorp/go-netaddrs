// Package netaddrs provides a function to get IP addresses given a
// 1. DNS name, OR
// 2. custom executable with optional args which
//  a. on success - exits 0 and prints whitespace delimited IP addresses to stdout.
//  b. on failure - exits with a non-zero code and/or optionally prints an error message of up to 1024 bytes to stderr.
package netaddrs

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
)

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// resolveDNS resolves the given DNS name and returns IP addresses
func resolveDNS(cfg string, l *log.Logger) ([]net.IPAddr, error) {
	ips, err := net.LookupIP(cfg)
	if err != nil {
		return nil, fmt.Errorf("Could not resolve DNS name: %s. error: %s", cfg, err)
	}

	var addrs []net.IPAddr
	for _, ip := range ips {
		addrs = append(addrs, net.IPAddr{IP: ip})
	}

	l.Println("Resolved DNS name", cfg, "to IP addresses", addrs)

	return addrs, nil
}

// execCmd returns IP addresses by running a custom executable with optional args specified in format
// "exec=<executable with optional args>", which
//  a. on success - exits 0 and prints whitespace delimited IP addresses to stdout.
//  b. on failure - exits with a non-zero code and/or optionally prints an error message of up to 1024 bytes to stderr.
func execCmd(cfg string, l *log.Logger) ([]net.IPAddr, error) {
	// exec=<executable arg1 arg2>
	executableWithArgs := strings.Split(cfg, "exec=")[1]
	commandWithArgs := strings.Fields(executableWithArgs)
	// validate length of commandWithArgs is >= 1 (command required, args optional)
	command, cmdArgs := commandWithArgs[0], commandWithArgs[1:]

	l.Println("Executing command: ", command, "Args list: ", cmdArgs)

	cmd := exec.Command(command, cmdArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("Error retrieving addresses on running the executable. Exit code: %d. Error message: %s", exitError.ExitCode(), errStr)
		}
		return nil, fmt.Errorf("Error retrieving addresses on running the executable. %s", err)
	}

	if len(errStr) > 0 {
		return nil, fmt.Errorf("Error retrieving addresses on running the executable. Details: %s", errStr)
	}

	execAddrs := strings.Fields(outStr)

	if len(execAddrs) < 1 {
		return nil, fmt.Errorf("Error retrieving addresses. Executable output: %s", outStr)
	}

	var addrs []net.IPAddr
	for _, addr := range execAddrs {
		addr = trimQuotes(addr)
		splitaddr := strings.Split(addr, "%")
		ipaddr := net.ParseIP(splitaddr[0])
		if ipaddr == nil {
			return nil, fmt.Errorf("Invalid IP address: %s.", splitaddr[0])
		}
		if len(splitaddr) == 2 {
			// ipv6 address
			addrs = append(addrs, net.IPAddr{IP: ipaddr, Zone: splitaddr[1]})
		} else {
			addrs = append(addrs, net.IPAddr{IP: ipaddr})
		}
	}

	l.Println("Addresses retrieved from the executable: ", addrs)

	return addrs, nil
}

// IPAddrs returns IP addresses given a
// 1. DNS name, OR
// 2. custom executable with optional args specified in format "exec=<executable with optional args>", which
//  a. on success - exits 0 and prints whitespace delimited IP addresses to stdout.
//  b. on failure - exits with a non-zero code and/or optionally prints an error message of up to 1024 bytes to stderr.
func IPAddrs(cfg string, l *log.Logger) ([]net.IPAddr, error) {

	if !strings.HasPrefix(cfg, "exec=") {
		return resolveDNS(cfg, l)
	}

	return execCmd(cfg, l)
}
