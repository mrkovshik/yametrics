package api

import "context"

// Server represents an interface for running a server.
type Server interface {
	// RunServer starts the server with the given context.
	// Parameters:
	// - ctx: the context to control server shutdown and other operations.
	RunServer(ctx context.Context)
}
