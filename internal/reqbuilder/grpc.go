package reqbuilder

import (
	"context"

	"github.com/mrkovshik/yametrics/internal/signature"
	"google.golang.org/grpc/metadata"
)

// GRPCContextBuilder is a builder for constructing gRPC contexts with metadata and HMAC-SHA256 signatures.
type GRPCContextBuilder struct {
	Ctx context.Context
	Err error
}

// NewGRPCContextBuilder creates a new instance of GRPCContextBuilder.
//
// Returns:
//   - *GRPCContextBuilder: A new GRPCContextBuilder instance with a background context.
func NewGRPCContextBuilder() *GRPCContextBuilder {
	ctx := context.Background()
	return &GRPCContextBuilder{
		Ctx: ctx,
	}
}

// AddMetaData adds metadata to the gRPC context.
//
// Parameters:
//   - key: The metadata key.
//   - value: The metadata value.
//
// Returns:
//   - *GRPCContextBuilder: The updated GRPCContextBuilder instance.
func (rb *GRPCContextBuilder) AddMetaData(key, value string) *GRPCContextBuilder {
	if rb.Err == nil {
		rb.Ctx = metadata.NewOutgoingContext(rb.Ctx, metadata.Pairs(key, value))
	}
	return rb
}

// Sign generates an HMAC-SHA256 signature for the message and adds it as metadata to the gRPC context.
//
// Parameters:
//   - key: The secret key for HMAC-SHA256 authentication.
//   - message: The message to be signed.
//
// Returns:
//   - *GRPCContextBuilder: The updated GRPCContextBuilder instance.
func (rb *GRPCContextBuilder) Sign(key string, message []byte) *GRPCContextBuilder {
	if key != "" && rb.Err == nil && len(message) != 0 {
		sigSrv := signature.NewSha256Sig(key, message)
		sig, err := sigSrv.Generate()
		if err != nil {
			rb.Err = err
			return rb
		}
		return rb.AddMetaData("HashSHA256", sig)
	}
	return rb
}
