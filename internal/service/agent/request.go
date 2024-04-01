package service

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type RequestBuilder struct {
	R   http.Request
	Err error
}

func NewRequestBuilder() *RequestBuilder {
	req, err := http.NewRequest(http.MethodGet, "", nil)
	return &RequestBuilder{*req, err}
}

func (rb *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
	rb.R.Header.Add(key, value)
	return rb
}

func (rb *RequestBuilder) SetMethod(method string) *RequestBuilder {
	if rb.Err == nil {
		rb.R.Method = method
	}
	return rb
}

func (rb *RequestBuilder) SetURL(rawURL string) *RequestBuilder {
	if rb.Err == nil {
		rb.R.URL, rb.Err = url.Parse(rawURL)
	}
	return rb
}

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