package rest

import (
	"encoding/json"
	"net/http"

	"github.com/mrkovshik/yametrics/internal/apperrors"
	"go.uber.org/zap"

	"github.com/mrkovshik/yametrics/internal/model"
)

// HandleUpdateMetricFromJSON handles HTTP requests to update a metric from JSON data.
func (s *Server) HandleUpdateMetricFromJSON(w http.ResponseWriter, r *http.Request) {
	var newMetrics model.Metrics
	ctx := r.Context()
	if err := newMetrics.MapMetricsFromReqJSON(r); err != nil {
		s.logger.Error("MapMetricsFromReqJSON", zap.Error(err))
		http.Error(w, apperrors.ErrInvalidRequestData.Error(), http.StatusBadRequest)
		return

	}

	if err := s.service.UpdateMetrics(ctx, []model.Metrics{newMetrics}); err != nil {
		s.logger.Error("UpdateMetrics", zap.Error(err))
		http.Error(w, "error w.Write", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	s.writeStatusWithMessage(w, http.StatusOK, "Gauge successfully updated")

}

// HandleUpdateMetricsFromJSON handles HTTP requests to update multiple metrics from JSON data.
func (s *Server) HandleUpdateMetricsFromJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var batch []model.Metrics
	if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
		s.logger.Error("Decode", zap.Error(err))
		http.Error(w, "Decode", http.StatusInternalServerError)
		return
	}
	if err := s.service.UpdateMetrics(ctx, batch); err != nil {
		s.logger.Error("UpdateMetrics", zap.Error(err))
		http.Error(w, "UpdateMetrics", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	s.writeStatusWithMessage(w, http.StatusOK, "Gauge successfully updated")
}

// HandleUpdateMetricFromURL handles HTTP requests to update a metric from URL parameters.
func (s *Server) HandleUpdateMetricFromURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newMetrics model.Metrics
	if err := newMetrics.MapMetricsFromReqURL(r); err != nil {
		s.logger.Error("MapMetricsFromReq", zap.Error(err))
		http.Error(w, apperrors.ErrInvalidRequestData.Error(), http.StatusBadRequest)
		return
	}
	if err := s.service.UpdateMetrics(ctx, []model.Metrics{newMetrics}); err != nil {
		s.logger.Error("UpdateMetrics", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	s.writeStatusWithMessage(w, http.StatusOK, "Gauge successfully updated")
}
