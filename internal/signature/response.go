package signature

import (
	"io"
	"net/http"
)

type CapturingResponseWriter struct {
	http.ResponseWriter
	buf []byte
}

func NewCapturingResponseWriter(w http.ResponseWriter) *CapturingResponseWriter {
	return &CapturingResponseWriter{
		ResponseWriter: w,
	}
}

func (rw *CapturingResponseWriter) Write(b []byte) (int, error) {
	rw.buf = append(rw.buf, b...)
	return rw.ResponseWriter.Write(b)
}

func (rw *CapturingResponseWriter) Body() []byte {
	return rw.buf
}
func (rw *CapturingResponseWriter) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(rw.buf)
	return int64(n), err
}
