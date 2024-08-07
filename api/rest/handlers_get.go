package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/mrkovshik/yametrics/internal/apperrors"
	"go.uber.org/zap"

	"github.com/mrkovshik/yametrics/internal/model"
)

// HandleGetMetricFromJSON handles HTTP requests to retrieve a metric using JSON data.
func (s *Server) HandleGetMetricFromJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newMetrics model.Metrics
	if err1 := json.NewDecoder(r.Body).Decode(&newMetrics); err1 != nil {
		s.logger.Error("Decode", zap.Error(err1))
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}
	metric, err2 := s.service.GetMetric(ctx, newMetrics)
	if err2 != nil {
		s.logger.Error("GetMetric", zap.Error(err2))
		http.Error(w, "GetMetric", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err3 := json.NewEncoder(w).Encode(metric); err3 != nil {
		s.logger.Error("Encode", zap.Error(err3))
		http.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleGetMetricFromURL handles HTTP requests to retrieve a metric using URL parameters.
func (s *Server) HandleGetMetricFromURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newMetrics model.Metrics
	if err := newMetrics.MapMetricsFromReqURL(r); err != nil {
		s.logger.Error("MapMetricsFromReq", zap.Error(err))
		http.Error(w, apperrors.ErrInvalidRequestData.Error(), http.StatusBadRequest)
		return
	}

	metric, err2 := s.service.GetMetric(ctx, newMetrics)
	if err2 != nil {
		s.logger.Error("s.storage.GetMetricByModel", zap.Error(err2))
		http.Error(w, "error getting value from server", http.StatusNotFound)
		return
	}

	var stringValue string
	switch metric.MType {
	case model.MetricTypeCounter:
		stringValue = strconv.FormatInt(*metric.Delta, 10)
	case model.MetricTypeGauge:
		stringValue = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	default:
		s.logger.Error("invalid metric type", zap.Error(errors.New("ErrInvalidMetricType")))
		http.Error(w, "error w.Write", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	s.writeStatusWithMessage(w, http.StatusOK, stringValue)
}

// HandleGetMetrics handles HTTP requests to retrieve all metrics.
func (s *Server) HandleGetMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "text/html")
	body, err := s.service.GetAllMetrics(ctx)
	if err != nil {
		s.logger.Error("s.storage.GetAllMetrics", zap.Error(err))
		http.Error(w, "s.storage.GetAllMetrics", http.StatusInternalServerError)
		return
	}
	s.writeStatusWithMessage(w, http.StatusOK, body)
}
