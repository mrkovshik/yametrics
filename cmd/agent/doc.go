// Package main provides the main entry point for running the metrics agent.
//
// The main function initializes the necessary components such as storage, metrics source, logging, and configuration.
// It starts the agent, which periodically polls metrics from a runtime source, stores them using a storage backend,
// and sends reports at specified intervals.
//
// Configuration:
// The configuration is loaded using the config package. It fetches configurations from environment variables and
// command-line flags. The agent's address, polling interval, reporting interval, authentication key, and rate limit
// are configured through these mechanisms.
//
// Running the Agent:
// The agent starts by initializing a map storage and runtime metrics source. It configures logging using zap, a
// structured logging library, and retrieves configuration settings using the GetConfigs function from the config package.
// Polling for metrics and sending reports are handled by separate goroutines started with Go's built-in concurrency support.
// The main function runs indefinitely, ensuring the agent continues to operate until interrupted.
//
// Example Usage:
// To run the metrics agent, configure environment variables or use command-line flags to set necessary parameters,
// such as the server address, polling interval, reporting interval, authentication key, and rate limit.
// Then run the binary built from this package.
//
// Usage:
//     go run main.go -a=localhost:8080 -p=2 -r=10 -k=your_auth_key -l=1

package main
