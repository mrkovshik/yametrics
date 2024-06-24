package api

import "context"

type Server interface {
	RunServer(ctx context.Context)
}
