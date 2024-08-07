// Package api defines interfaces and utilities for server-related operations.
package api

import (
	"context"
)

// Server represents an interface for running a server.
type Server interface {
	// RunServer starts the server with the given context.
	RunServer(ctx context.Context) error
}
