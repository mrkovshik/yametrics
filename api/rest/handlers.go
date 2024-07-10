package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/mrkovshik/yametrics/internal/model"
)

var errInvalidRequestData = errors.New("invalid request data")

// HandleUpdateMetricFromJSON handles HTTP requests to update a metric from JSON data.
func (s *Server) HandleUpdateMetricFromJSON(w http.ResponseWriter, r *http.Request) {
	var newMetrics model.Metrics
	ctx := r.Context()
	if err := newMetrics.MapMetricsFromReqJSON(r); err != nil {
		s.logger.Error("MapMetricsFromReqJSON", zap.Error(err))
		http.Error(w, errInvalidRequestData.Error(), http.StatusBadRequest)
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
		http.Error(w, errInvalidRequestData.Error(), http.StatusBadRequest)
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
		http.Error(w, errInvalidRequestData.Error(), http.StatusBadRequest)
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
		stringValue = fmt.Sprint(*metric.Delta)
	case model.MetricTypeGauge:
		stringValue = fmt.Sprint(*metric.Value)
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

func (s *Server) writeStatusWithMessage(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	if _, err := w.Write([]byte(msg)); err != nil {
		s.logger.Error("w.Write:", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
