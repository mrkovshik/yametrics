package compress

import (
	"compress/gzip"
	"io"
	"net/http"
)

type GzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewGzipWriter(w http.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}
func (c *GzipWriter) Header() http.Header {
	return c.w.Header()
}

func (c *GzipWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *GzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *GzipWriter) Close() error {
	return c.zw.Close()
}

type GzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

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

func (c GzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *GzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
