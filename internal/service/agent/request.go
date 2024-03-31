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
	return &RequestBuilder{http.Request{}, nil}
}

func (rb RequestBuilder) WithHeaders(headers map[string]string) *RequestBuilder {
	if rb.err == nil {
		for header, value := range headers {
			rb.r.Header.Add(header, value)
		}
	}
	return &rb
}

func (rb RequestBuilder) SetMethod(method string) *RequestBuilder {
	if rb.err == nil {
		rb.r.Method = method
	}
	return &rb
}

func (rb RequestBuilder) SetURL(rawUrl string) *RequestBuilder {
	if rb.err == nil {
		rb.r.URL, rb.err = url.Parse(rawUrl)
	}
	return &rb
}

func (rb RequestBuilder) AddJSONBody(data any) *RequestBuilder {
	if rb.err == nil {
		buf := bytes.Buffer{}
		rb.err = json.NewEncoder(&buf).Encode(data)
		if rb.err == nil {
			rb.r.Body = io.NopCloser(&buf)
			rb.WithHeaders(map[string]string{"Content-Type": "application/json"})
		}
	}
	return &rb
}

func (rb RequestBuilder) Compress() *RequestBuilder {
	if rb.err == nil {
		compressedBody, err := gzip.NewReader(rb.r.Body)
		if err == nil {
			rb.r.Body = compressedBody
			rb.WithHeaders(map[string]string{"Content-Encoding": "gzip"})
		}
	}
	return &rb
}
