package service

import (
	"context"
	"time"
)

func (s *Server) DumpMetrics(ctx context.Context) {
	for {
		time.Sleep(time.Second * time.Duration(s.config.StoreInterval))
		if err := s.storage.StoreMetrics(ctx, s.config.StoreFilePath); err != nil {
			s.logger.Error("StoreMetrics", err)
		}

	}
}

func (s *Server) StoreMetrics(ctx context.Context, path string) error {
	return s.storage.StoreMetrics(ctx, path)
}

func (s *Server) RestoreMetrics(ctx context.Context, path string) error {
	return s.storage.RestoreMetrics(ctx, path)
}
