// Package netaddrs provides a function to get IP addresses from a DNS request or
// by executing a binary.
package netaddrs

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
)

// IPAddrs looks up and returns IP addresses using the method described by cfg.
//
// If cfg is a DNS name IP addresses are looked up by querying the default
// DNS resolver for A and AAAA records associated with the DNS name.
//
// If cfg has an exec= prefix, IP addresses are looked up by executing the command
// after exec=. The command may include optional arguments. Command arguments
// must be space separated (spaces in argument values can not be escaped).
// The command may output IPv4 or IPv6 addresses, and IPv6 addresses can
// optionally include a zone index.
// The executable must follow these rules:
//
//    on success - exit 0 and print whitespace delimited IP addresses to stdout.
//    on failure - exits with a non-zero code, and should print an error message
//                 of up to 1024 bytes to stderr.
func IPAddrs(cfg string, l *log.Logger) ([]net.IPAddr, error) {
	if !strings.HasPrefix(cfg, "exec=") {
		return resolveDNS(cfg, l)
	}

	ips, err := execCmd(cfg, l)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve IP addresses from executable: %w", err)
	}
	return ips, nil
}

// resolveDNS resolves the given DNS name and returns IP addresses
func resolveDNS(cfg string, l *log.Logger) ([]net.IPAddr, error) {
	ips, err := net.LookupIP(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve DNS name: %s: %s", cfg, err)
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

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("executable failed with exit code %d: %s", exitError.ExitCode(), strings.TrimSpace(stderr.String()))
		}
		return nil, err
	}

	execAddrs := strings.Fields(stdout.String())
	if len(execAddrs) < 1 {
		return nil, fmt.Errorf("executable returned no output to stdout")
	}

	var addrs []net.IPAddr
	for _, addr := range execAddrs {
		addr = trimQuotes(addr)
		splitaddr := strings.Split(addr, "%")
		ipaddr := net.ParseIP(splitaddr[0])
		if ipaddr == nil {
			return nil, fmt.Errorf("executable returned invalid IP address: %s", splitaddr[0])
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

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}
