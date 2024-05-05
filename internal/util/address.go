package util

import (
	"net"
	"strings"
)

func ValidateAddress(addr string) bool {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return false
	}
	host := parts[0]
	port := parts[1]
	_, err := net.LookupPort("tcp", port)
	if err != nil {
		return false
	}
	ip := net.ParseIP(host)
	if ip == nil {
		if _, err := net.LookupIP(host); err != nil {
			return false
		}
	}
	return true
}
