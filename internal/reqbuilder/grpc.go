package reqbuilder

import (
	"context"

	"github.com/mrkovshik/yametrics/internal/signature"
	"google.golang.org/grpc/metadata"
)

type GRPCContextBuilder struct {
	Ctx context.Context
	Err error
}

func NewGRPCContextBuilder() *GRPCContextBuilder {
	ctx := context.Background()
	return &GRPCContextBuilder{
		Ctx: ctx,
	}
}

func (rb *GRPCContextBuilder) AddMetaData(key, value string) *GRPCContextBuilder {
	if rb.Err == nil {
		rb.Ctx = metadata.NewOutgoingContext(rb.Ctx, metadata.Pairs(key, value))
	}
	return rb
}

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
