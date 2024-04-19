package signature

import (
	"io"
	"net/http"
)

type CapturingResponseWriter struct {
	http.ResponseWriter
	statCode int
	buf      []byte
}

func NewCapturingResponseWriter(w http.ResponseWriter) *CapturingResponseWriter {
	return &CapturingResponseWriter{
		ResponseWriter: w,
	}
}

func (rw *CapturingResponseWriter) WriteHeader(statusCode int) {
	rw.statCode = statusCode
}

func (rw *CapturingResponseWriter) Write(b []byte) (int, error) {
	rw.buf = append(rw.buf, b...)
	return 0, nil
}

func (rw *CapturingResponseWriter) Body() []byte {
	return rw.buf
}

func (rw *CapturingResponseWriter) Code() int {
	return rw.statCode
}

func (rw *CapturingResponseWriter) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(rw.buf)
	return int64(n), err
}
