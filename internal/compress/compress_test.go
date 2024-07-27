package compress

import (
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipWriter_WriteHeader(t *testing.T) {
	// Create a mock ResponseWriter
	w := httptest.NewRecorder()

	// Create a new gzip writer
	zw := gzip.NewWriter(w)

	// Create a GzipWriter instance
	c := &GzipWriter{
		w:  w,
		zw: zw,
	}

	type args struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{http.StatusOK}, "gzip"},
		{"2", args{http.StatusBadRequest}, ""},
	}
	for _, tt := range tests {
		c.w.Header().Set("Content-Encoding", "")
		t.Run(tt.name, func(t *testing.T) {
			c.WriteHeader(tt.args.statusCode)
			assert.Equal(t, tt.want, c.w.Header().Get("Content-Encoding"))
		})
	}
}
