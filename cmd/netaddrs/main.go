package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
		w = ioutil.Discard
	}
	l := log.New(w, "netaddrs: ", 0)
	addresses, err := netaddrs.IPAddrs(args[1], l)

	if err != nil {
		l.Fatal(err)
	}

	var outputAddresses []string
	for _, address := range addresses {
		outputAddresses = append(outputAddresses, address.IP.String())
	}

	fmt.Println(strings.Join(outputAddresses, " "))
}
