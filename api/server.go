package api

// Server represents an interface for running a server.
type Server interface {
	// RunServer starts the server with the given context.
	RunServer()
}
