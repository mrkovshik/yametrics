package service

import "time"

func (s *Server) DumpMetrics() {
	for {
		time.Sleep(time.Second * time.Duration(s.config.StoreInterval))
		if err := s.storage.StoreMetrics(s.config.StoreFilePath); err != nil {
			s.logger.Error("StoreMetrics", err)
		}

	}
}

func (s *Server) StoreMetrics(path string) error {
	return s.storage.StoreMetrics(path)
}

func (s *Server) RestoreMetrics(path string) error {
	return s.storage.RestoreMetrics(path)
}
