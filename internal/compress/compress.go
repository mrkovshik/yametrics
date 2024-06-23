// Package compress is intended to decrease the size of the requests and responses.
package compress

import (
	"compress/gzip"
	"io"
	"net/http"
)

// GzipWriter wraps an http.ResponseWriter to provide gzip compression.
type GzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewGzipWriter creates a new GzipWriter that wraps the provided http.ResponseWriter.
func NewGzipWriter(w http.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header returns the header map that will be sent by WriteHeader.
func (c *GzipWriter) Header() http.Header {
	return c.w.Header()
}

// Write compresses the data and writes it to the wrapped http.ResponseWriter.
func (c *GzipWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader sends an HTTP response header with the provided status code.
// If the status code indicates a successful response (< 300), it sets the
// "Content-Encoding" header to "gzip".
func (c *GzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close closes the gzip.Writer, flushing any unwritten data to the underlying
// http.ResponseWriter.
func (c *GzipWriter) Close() error {
	return c.zw.Close()
}

// GzipReader wraps an io.ReadCloser to provide gzip decompression.
type GzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewGzipReader creates a new GzipReader that wraps the provided io.ReadCloser.
// It returns an error if the gzip reader could not be created.
func NewGzipReader(r io.ReadCloser) (*GzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &GzipReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read reads and decompresses data from the wrapped io.ReadCloser.
func (c GzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closes both the underlying io.ReadCloser and the gzip.Reader.
func (c *GzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
