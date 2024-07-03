// Package api defines interfaces and utilities for server-related operations.
package api

import (
	"os"
)

// Server represents an interface for running a server.
type Server interface {
	// RunServer starts the server with the given context.
	RunServer(stop chan os.Signal) error
}
