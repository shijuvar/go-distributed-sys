package network

import (
	// stdlib
	"errors"
	"net"
	"os"
	"strings"
)

// ErrIFNotFound is returned if we can't find an available network interface.
var ErrIFNotFound = errors.New("network interface not found")

// HostIP tries to return the IPv4 string from the primary network interface.
func HostIP() (string, error) {
	_, ip, err := getPrimaryIPv4Interface()
	if err != nil {
		// fallback to hostname
		hostName, err := os.Hostname()
		if err != nil {
			return "", err
		}
		ips, err := net.LookupIP(hostName)
		if err != nil {
			return "", err
		}
		for _, ip2 := range ips {
			if ip2.To4() != nil {
				ip = &ip2
				break
			}
		}
	}
	if ip == nil {
		return "", ErrIFNotFound
	}
	return ip.String(), nil
}

func getPrimaryIPv4Interface() (*net.Interface, *net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}
	for _, i := range interfaces {
		if i.Flags&net.FlagUp != net.FlagUp || i.Flags&net.FlagLoopback == net.FlagLoopback {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			tokens := strings.Split(a.String(), "/") // ip/netmask
			ip := net.ParseIP(tokens[0])
			if ip.To4() != nil {
				return &i, &ip, nil
			}
		}
	}
	return nil, nil, ErrIFNotFound
}
