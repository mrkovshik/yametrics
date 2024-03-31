package service

import "time"

func (s *Server) DumpMetrics() {
	for {
		if err := s.Storage.DumpMetrics(s.Config.StoreFilePath); err != nil {
			s.Logger.Error("DumpMetrics", err)
		}
		time.Sleep(time.Second * time.Duration(s.Config.StoreInterval))
	}
}
