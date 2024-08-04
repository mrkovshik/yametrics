// Package grpc provides a gRPC server implementation for the yametrics application.
// It includes functionalities to start, run, and gracefully stop the server,
// as well as handling gRPC requests and interceptors for logging and authentication.
//
// The Server type encapsulates the gRPC server, the service implementing business logic,
// configuration settings, and logging.
// It is used to initialize and run the gRPC server,
// and the RunServer method handles starting the server and listening for shutdown signals.
package grpc
