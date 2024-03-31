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
	r   http.Request
	err error
}

func NewRequestBuilder() *RequestBuilder {
	req, err := http.NewRequest(http.MethodGet, "", nil)
	return &RequestBuilder{*req, err}
}

func (rb *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
	rb.r.Header.Add(key, value)
	return rb
}

func (rb *RequestBuilder) SetMethod(method string) *RequestBuilder {
	if rb.err == nil {
		rb.r.Method = method
	}
	return rb
}

func (rb *RequestBuilder) SetURL(rawUrl string) *RequestBuilder {
	if rb.err == nil {
		rb.r.URL, rb.err = url.Parse(rawUrl)
	}
	return rb
}

func (rb *RequestBuilder) AddJSONBody(data any) *RequestBuilder {
	if rb.err == nil {
		buf := bytes.Buffer{}
		rb.err = json.NewEncoder(&buf).Encode(data)
		if rb.err == nil {
			rb.r.Body = io.NopCloser(&buf)
			rb.WithHeader("Content-Type", "application/json")
		}
	}
	return rb
}

func (rb *RequestBuilder) Compress() *RequestBuilder {
	if rb.err == nil {
		var compressedBody bytes.Buffer
		gzipWriter := gzip.NewWriter(&compressedBody)
		_, err := io.Copy(gzipWriter, rb.r.Body)
		if err != nil {
			return &RequestBuilder{err: err}
		}
		err = gzipWriter.Close()
		if err != nil {
			return &RequestBuilder{err: err}
		}
		rb.r.Body = io.NopCloser(&compressedBody)
		rb.r.ContentLength = int64(compressedBody.Len())
		rb.WithHeader("Content-Encoding", "gzip")
	}
	return rb
}
