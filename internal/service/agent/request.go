package service

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	rsa2 "github.com/mrkovshik/yametrics/internal/rsa"
	"github.com/mrkovshik/yametrics/internal/signature"
)

// RequestBuilder helps in constructing and modifying HTTP requests.
type RequestBuilder struct {
	R   http.Request // The HTTP request being built.
	Err error        // Any error encountered during the building process.
}

// NewRequestBuilder initializes a new RequestBuilder with a default GET request.
func NewRequestBuilder() *RequestBuilder {
	req, err := http.NewRequest(http.MethodGet, "", nil)
	return &RequestBuilder{*req, err}
}

// WithHeader adds a header to the HTTP request.
func (rb *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
	rb.R.Header.Add(key, value)
	return rb
}

// SetMethod sets the HTTP method for the request.
func (rb *RequestBuilder) SetMethod(method string) *RequestBuilder {
	if rb.Err == nil {
		rb.R.Method = method
	}
	return rb
}

// SetURL sets the URL for the HTTP request.
func (rb *RequestBuilder) SetURL(rawURL string) *RequestBuilder {
	if rb.Err == nil {
		rb.R.URL, rb.Err = url.Parse(rawURL)
	}
	return rb
}

// AddJSONBody encodes data as JSON and sets it as the body of the request.
func (rb *RequestBuilder) AddJSONBody(data any) *RequestBuilder {
	if rb.Err == nil && data != nil {
		buf := bytes.Buffer{}
		rb.Err = json.NewEncoder(&buf).Encode(data)
		if rb.Err == nil {
			rb.R.Body = io.NopCloser(&buf)
			rb.WithHeader("Content-Type", "application/json")
		}
	}
	return rb
}

// Sign generates a SHA-256 signature for the request body and adds it as a header.
func (rb *RequestBuilder) Sign(key string) *RequestBuilder {
	var body []byte
	if key != "" && rb.Err == nil && rb.R.Body != nil {
		body, rb.Err = io.ReadAll(rb.R.Body)
		rb.R.Body = io.NopCloser(bytes.NewBuffer(body))
		if rb.Err == nil {
			sigSrv := signature.NewSha256Sig(key, body)
			sig, err := sigSrv.Generate()
			if err != nil {
				rb.Err = err
				return rb
			}
			rb.WithHeader("HashSHA256", sig)
		}
	}
	return rb
}

// Compress compresses the request body using gzip and sets the appropriate headers.
func (rb *RequestBuilder) Compress() *RequestBuilder {
	if rb.Err == nil {
		var compressedBody bytes.Buffer
		gzipWriter := gzip.NewWriter(&compressedBody)
		_, err := io.Copy(gzipWriter, rb.R.Body)
		if err != nil {
			return &RequestBuilder{Err: err}
		}
		err = gzipWriter.Close()
		if err != nil {
			return &RequestBuilder{Err: err}
		}
		rb.R.Body = io.NopCloser(&compressedBody)
		rb.R.ContentLength = int64(compressedBody.Len())
		rb.WithHeader("Content-Encoding", "gzip")
	}
	return rb
}

func (rb *RequestBuilder) EncryptRSA(pemFilePath string) *RequestBuilder {
	var body []byte
	if pemFilePath != "" && rb.Err == nil && rb.R.Body != nil {
		// Read the request body
		body, rb.Err = io.ReadAll(rb.R.Body)

		if rb.Err == nil {
			// Read the PEM file
			publicKeyPem, err := rsa2.ReadPEMFile(pemFilePath)
			if err != nil {
				rb.Err = err
				return rb
			}

			// Encrypt the body using RSA
			encryptedBody, err := rsa2.Encrypt(publicKeyPem, body)
			if err != nil {
				rb.Err = err
				return rb
			}
			rb.R.Body = io.NopCloser(bytes.NewBufferString(encryptedBody))
		}
	}
	return rb
}
