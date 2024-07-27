// Package rest provides the implementation for the HTTP server handling the yametrics service.
// It includes various HTTP handlers and middleware for managing metrics, including updating metrics,
// retrieving metrics, and performing server/database health checks.
//
// This package uses the go-chi/chi router for routing and provides middleware functionalities for:
//
// - Logging: Logs incoming HTTP requests and their corresponding responses.
// - Gzip Compression: Handles gzip compression for request and response bodies.
// - Authentication: Authenticates incoming requests using HMAC-SHA256 signatures.
// - Response Signing: Signs outgoing response bodies using HMAC-SHA256 signatures if a signing key is configured.
//
// ## Components
//
// The primary components of this package are:
//
// - Server: Represents the server configuration and dependencies.
// - NewServer: Creates a new Server instance with the provided service, configuration, and logger.
// - ConfigureRouter: configures routes and middleware.
// - RunServer: Starts the HTTP server with the configured router.
//
// ## Routes
//
// The following routes are handled by the server:
//
// - POST /update/: Updates a single metric from JSON data.
// - POST /update/{type}/{name}/{value}: Updates a single metric from URL parameters.
// - POST /updates/: Updates multiple metrics from JSON data.
// - POST /value/: Retrieves a single metric using JSON data.
// - GET /value/{type}/{name}: Retrieves a single metric using URL parameters.
// - GET /ping: Checks the health of the server/database.
// - GET /: Retrieves all metrics.
//
// ## Middleware
//
// Middleware functionalities include:
//
// - WithLogging: Logs incoming HTTP requests and their responses.
// - GzipHandle: Manages gzip compression for request and response bodies.
// - Authenticate: Verifies the integrity of incoming requests using HMAC-SHA256 signatures.
// - SignResponse: Signs outgoing response bodies using HMAC-SHA256 signatures if a signing key is configured.
package rest
