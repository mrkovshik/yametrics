package rest

import (
	"bytes"
	"crypto/hmac"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mrkovshik/yametrics/internal/compress"
	"github.com/mrkovshik/yametrics/internal/logger"
	"github.com/mrkovshik/yametrics/internal/signature"
)

// WithLogging wraps an http.Handler with logging functionality.
// It logs incoming HTTP requests and their corresponding responses.
func (s *Server) WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &logger.ResponseData{
			Status: 0,
			Size:   0,
		}
		lw := logger.LoggingResponseWriter{
			ResponseWriter: w,
			ResponseData:   responseData,
		}
		h.ServeHTTP(&lw, r)
		duration := time.Since(start)
		s.logger.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.Status,
			"duration", duration,
			"size", responseData.Size,
		)
	}
	return http.HandlerFunc(logFn)
}

// GzipHandle returns an http.Handler that handles gzip compression for request and response bodies.
// It checks the request headers for gzip encoding and wraps the response writer with gzip compression if supported.
func (s *Server) GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var isEncodingSupported = false

		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := compress.NewGzipReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = gz
			defer gz.Close() //nolint:all
		}

		acceptValues := r.Header.Values("Accept-Encoding")
		for _, value := range acceptValues {
			if strings.Contains(value, "gzip") {
				isEncodingSupported = true
				break
			}
		}

		if !isEncodingSupported {
			next.ServeHTTP(w, r)
			return
		}

		cw := compress.NewGzipWriter(w)

		defer cw.Close() //nolint:all

		next.ServeHTTP(cw, r)
	})
}

// Authenticate returns an http.Handler that authenticates incoming requests using HMAC-SHA256 signatures.
// It verifies the integrity of the request body against the provided signature.
func (s *Server) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientSig := r.Header.Get(`HashSHA256`)
		if clientSig != "" && r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			sigSrv := signature.NewSha256Sig(s.config.Key, body)
			sig, err := sigSrv.Generate()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !hmac.Equal([]byte(clientSig), []byte(sig)) {
				http.Error(w, "invalid signature", http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// SignResponse returns an http.Handler that signs outgoing response bodies using HMAC-SHA256 signatures.
// If a signing key is configured, it computes the signature of the response body and sets the HashSHA256 header.
func (s *Server) SignResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.config.Key == "" {
			next.ServeHTTP(w, r)
			return
		}
		rw := signature.NewCapturingResponseWriter(w)
		next.ServeHTTP(rw, r)
		if len(rw.Body()) != 0 {
			sigSrv := signature.NewSha256Sig(s.config.Key, rw.Body())
			sig, err := sigSrv.Generate()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("HashSHA256", sig)
			w.WriteHeader(rw.Code())
			_, err = w.Write(rw.Body())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})
}
