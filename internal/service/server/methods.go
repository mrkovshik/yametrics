package service

import (
	"context"
	"time"
)

// DumpMetrics periodically stores metrics to a file based on configuration.
func (s *Server) DumpMetrics(ctx context.Context) {
	for {
		time.Sleep(time.Second * time.Duration(s.config.StoreInterval))
		if err := s.storage.StoreMetrics(ctx, s.config.StoreFilePath); err != nil {
			s.logger.Error("StoreMetrics", err)
		}
	}
}

// StoreMetrics delegates storing metrics to the storage interface.
func (s *Server) StoreMetrics(ctx context.Context, path string) error {
	return s.storage.StoreMetrics(ctx, path)
}

// RestoreMetrics delegates restoring metrics from a file to the storage interface.
func (s *Server) RestoreMetrics(ctx context.Context, path string) error {
	return s.storage.RestoreMetrics(ctx, path)
}
