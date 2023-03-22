// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	netaddrs "github.com/hashicorp/go-netaddrs"
)

func main() {
	var quiet bool
	var help bool
	flag.BoolVar(&quiet, "q", false, "no verbose output")
	flag.BoolVar(&help, "h", false, "print help")
	flag.Parse()

	args := flag.Args()
	if help || len(args) == 0 || args[0] != "ip" {
		fmt.Println("Usage: netaddrs ip \"DNS name\" or \"exec=<executable with optional args>\"")
		os.Exit(0)
	}

	var w io.Writer = os.Stderr
	if quiet {
		w = io.Discard
	}
	l := log.New(w, "netaddrs: ", 0)
	addresses, err := netaddrs.IPAddrs(context.Background(), args[1], logger{l: l})
	if err != nil {
		l.Fatal(err)
	}

	var outputAddresses []string
	for _, address := range addresses {
		outputAddress := address.IP.String()
		if address.Zone != "" {
			outputAddress += "%" + address.Zone
		}
		outputAddresses = append(outputAddresses, outputAddress)
	}

	fmt.Println(strings.Join(outputAddresses, " "))
}

type logger struct {
	l *log.Logger
}

func (l logger) Debug(msg string, args ...interface{}) {
	l.l.Print(msg)
	l.l.Println(args...)
}
