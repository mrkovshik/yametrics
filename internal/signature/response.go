// Package signature provides utilities for signing the request.
package signature

import (
	"io"
	"net/http"
)

// CapturingResponseWriter is a wrapper around http.ResponseWriter that captures response details.
type CapturingResponseWriter struct {
	http.ResponseWriter        // Embedding http.ResponseWriter to delegate standard HTTP response methods
	statCode            int    // Status code of the HTTP response
	buf                 []byte // Buffer to capture response body
}

// NewCapturingResponseWriter creates a new CapturingResponseWriter instance.
func NewCapturingResponseWriter(w http.ResponseWriter) *CapturingResponseWriter {
	return &CapturingResponseWriter{
		ResponseWriter: w,
	}
}

// WriteHeader captures the HTTP status code when it is set.
func (rw *CapturingResponseWriter) WriteHeader(statusCode int) {
	rw.statCode = statusCode
}

// Write captures the response body data.
// It appends the provided bytes to the internal buffer.
func (rw *CapturingResponseWriter) Write(b []byte) (int, error) {
	rw.buf = append(rw.buf, b...)
	return len(b), nil
}

// Body returns the captured response body.
func (rw *CapturingResponseWriter) Body() []byte {
	return rw.buf
}

// Code returns the captured HTTP status code.
func (rw *CapturingResponseWriter) Code() int {
	return rw.statCode
}

// WriteTo writes the captured response body to the given io.Writer.
func (rw *CapturingResponseWriter) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(rw.buf)
	return int64(n), err
}
