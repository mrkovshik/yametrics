package rest

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// HandlePing handles HTTP requests to ping the server/database.
func (s *Server) HandlePing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if s.config.DBEnable {
		newCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := s.service.Ping(newCtx); err != nil {
			s.logger.Error("PingContext", zap.Error(err))
			http.Error(w, "data base is not responding", http.StatusInternalServerError)
			return
		}
		s.writeStatusWithMessage(w, http.StatusOK, "database is alive")
	}
	s.writeStatusWithMessage(w, http.StatusInternalServerError, "DB is unable")
}
