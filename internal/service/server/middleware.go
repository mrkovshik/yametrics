package service

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
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
