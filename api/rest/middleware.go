package rest

import (
	"bytes"
	"crypto/hmac"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mrkovshik/yametrics/internal/compress"
	"github.com/mrkovshik/yametrics/internal/logger"
	rsa2 "github.com/mrkovshik/yametrics/internal/rsa"
	"github.com/mrkovshik/yametrics/internal/signature"
)

// WithLogging wraps an http.Handler with logging functionality.
// It logs incoming HTTP requests and their corresponding responses.
//
// Parameters:
//   - h: The http.Handler to be wrapped with logging.
//
// Returns:
//   - http.Handler: The wrapped http.Handler with logging.
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
		defer s.logger.Infoln(
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
//
// Parameters:
//   - next: The next http.Handler to be called.
//
// Returns:
//   - http.Handler: The wrapped http.Handler with gzip compression support.
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
			defer gz.Close() // nolint:all
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

		defer cw.Close() // nolint:all

		next.ServeHTTP(cw, r)
	})
}

// Authenticate returns an http.Handler that authenticates incoming requests using HMAC-SHA256 signatures.
// It verifies the integrity of the request body against the provided signature.
//
// Parameters:
//   - next: The next http.Handler to be called.
//
// Returns:
//   - http.Handler: The wrapped http.Handler with authentication.
func (s *Server) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientSig := r.Header.Get(`HashSHA256`)
		if clientSig != "" && r.Body != nil {
			body, err := io.ReadAll(r.Body)
			defer r.Body.Close() // nolint:all
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
//
// Parameters:
//   - next: The next http.Handler to be called.
//
// Returns:
//   - http.Handler: The wrapped http.Handler with response signing.
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

// DecryptRequest returns an http.Handler that decrypts incoming requests using RSA.
// If a decryption key is configured, it decrypts the request body and replaces the original body.
//
// Parameters:
//   - next: The next http.Handler to be called.
//
// Returns:
//   - http.Handler: The wrapped http.Handler with request decryption.
func (s *Server) DecryptRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.config.CryptoKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// Read the PEM file
		privateKeyPem, err := rsa2.ReadPEMFile(s.config.CryptoKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Decrypt the body using RSA
		plaintext, err := rsa2.Decrypt(privateKeyPem, body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(plaintext))

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// VerifySubnet returns an http.Handler that verifies the client IP address against a trusted subnet.
// If the client IP is not within the trusted subnet, it returns a forbidden error.
//
// Parameters:
//   - next: The next http.Handler to be called.
//
// Returns:
//   - http.Handler: The wrapped http.Handler with subnet verification.
func (s *Server) VerifySubnet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.config.TrustedSubnet == "" {
			next.ServeHTTP(w, r)
			return
		}
		clientIP := r.Header.Get(`X-Real-IP`)
		ip := net.ParseIP(clientIP)
		if ip == nil {
			http.Error(w, "Invalid IP address", http.StatusBadRequest)
			return
		}

		_, trustedNet, err := net.ParseCIDR(s.config.TrustedSubnet)
		if err != nil {
			http.Error(w, "Invalid trusted subnet", http.StatusInternalServerError)
			return
		}

		if !trustedNet.Contains(ip) {
			http.Error(w, "Forbidden: IP not in trusted subnet", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
