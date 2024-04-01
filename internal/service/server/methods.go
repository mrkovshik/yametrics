package service

import "time"

func (s *Server) DumpMetrics() {
	for {
		time.Sleep(time.Second * time.Duration(s.Config.StoreInterval))
		if err := s.Storage.StoreMetrics(s.Config.StoreFilePath); err != nil {
			s.Logger.Error("StoreMetrics", err)
		}

	}
}
