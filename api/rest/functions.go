package rest

import (
	"net/http"

	"go.uber.org/zap"
)

func (s *Server) writeStatusWithMessage(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	if _, err := w.Write([]byte(msg)); err != nil {
		s.logger.Error("w.Write:", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
