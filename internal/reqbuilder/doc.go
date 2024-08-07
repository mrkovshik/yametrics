// Package reqbuilder provides utilities for building gRPC and HTTP requests with metadata and signatures.
// It is designed to simplify the process of constructing requests, adding headers, encoding bodies,
// compressing data, and generating HMAC-SHA256 signatures and RSA encryption.
//
// The GRPCContextBuilder type is used for building gRPC contexts with metadata and HMAC-SHA256 signatures.
// The HTTPRequestBuilder type is used for constructing and modifying HTTP requests.
//
// Example usage for gRPC:
//
//	builder := reqbuilder.NewGRPCContextBuilder().
//	    AddMetaData("key", "value").
//	    Sign("secret_key", []byte("message"))
//
// Example usage for HTTP:
//
//	builder := reqbuilder.NewHTTPRequestBuilder().
//	    SetMethod("POST").
//	    SetURL("https://example.com").
//	    AddJSONBody(data).
//	    Sign("secret_key").
//	    Compress()
package reqbuilder
