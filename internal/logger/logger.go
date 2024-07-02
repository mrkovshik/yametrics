// Package logger provides middleware components for logging HTTP responses.
package logger

import "net/http"

// ResponseData holds information about the HTTP response status and size.
type ResponseData struct {
	Status int // HTTP status code of the response
	Size   int // Size of the response body in bytes
}

// LoggingResponseWriter wraps an http.ResponseWriter to capture response status and size.
type LoggingResponseWriter struct {
	http.ResponseWriter
	ResponseData *ResponseData
}

// Write writes the data to the connection as part of an HTTP reply and captures the size of the response.
func (r *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.ResponseData.Size += size // capture the size of the response
	return size, err
}

// WriteHeader sends an HTTP response header with the provided status code and captures the status.
func (r *LoggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.ResponseData.Status = statusCode // capture the status code
}
