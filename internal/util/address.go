package util

import (
	"net"
	"strings"
)

// ValidateAddress checks if the given address is valid.
// An address is considered valid if it contains a valid host and port.
// Parameters:
// - addr: the address to be validated.
// Returns:
// - true if the address is valid, false otherwise.
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
