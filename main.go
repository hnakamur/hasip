package main

import (
	"fmt"
	"net"
	"os"

	"github.com/alecthomas/kingpin"
)

var (
	intfName = kingpin.Flag("interface", "Name of ethernet interface.").Short('i').Required().String()
	addr     = kingpin.Flag("address", "IP address.").Short('a').Required().String()
)

func main() {
	kingpin.Parse()
	ok, err := run(*intfName, *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if !ok {
		os.Exit(1)
	}
	os.Exit(0)
}

func run(intfName, addr string) (bool, error) {
	intf, err := net.InterfaceByName(intfName)
	if err != nil {
		return false, fmt.Errorf("interface not found: %v", err)
	}

	ip := net.ParseIP(addr)
	ok, err := hasAddr(intf, ip)
	if err != nil {
		return false, fmt.Errorf("failed to check having ip: %v", err)
	}
	return ok, nil
}

func hasAddr(intf *net.Interface, ip net.IP) (bool, error) {
	addrs, err := intf.Addrs()
	if err != nil {
		return false, err
	}
	for _, addr := range addrs {
		i, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return false, err
		}
		if i.Equal(ip) {
			return true, nil
		}
	}
	return false, nil
}
