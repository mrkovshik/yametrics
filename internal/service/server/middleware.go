package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/mrkovshik/yametrics/internal/compress"
	"github.com/mrkovshik/yametrics/internal/logger"
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
		s.Logger.Infoln(
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
		//if cType := r.Header.Get("Content-Type"); cType != "application/json" && cType != "text/html" {
		//	next.ServeHTTP(w, r)
		//	return
		//}
		var (
			isEncodingSupported = false
		)

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
