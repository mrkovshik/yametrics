// Package service methods for handling various HTTP endpoints related to metrics and database operations.
package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mrkovshik/yametrics/internal/model"
	"go.uber.org/zap"
)

var errInvalidRequestData = errors.New("invalid request data")

// UpdateMetricFromJSON handles HTTP requests to update a metric from JSON data.
func (s *restAPIServer) UpdateMetricFromJSON(ctx context.Context) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var newMetrics model.Metrics
		if err := newMetrics.MapMetricsFromReqJSON(req); err != nil {
			s.logger.Error("MapMetricsFromReqJSON", zap.Error(err))
			http.Error(res, errInvalidRequestData.Error(), http.StatusBadRequest)
			return
		}

		if err := s.service.UpdateMetrics(ctx, []model.Metrics{newMetrics}); err != nil {
			s.logger.Error("UpdateMetrics", zap.Error(err))
			http.Error(res, "error res.Write", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write([]byte("Gauge successfully updated")); err != nil {
			s.logger.Error("res.Write", zap.Error(err))
			http.Error(res, "error res.Write", http.StatusInternalServerError)
			return
		}
	}
}

// UpdateMetricsFromJSON handles HTTP requests to update multiple metrics from JSON data.
func (s *restAPIServer) UpdateMetricsFromJSON(ctx context.Context) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var batch []model.Metrics
		if err := json.NewDecoder(req.Body).Decode(&batch); err != nil {
			s.logger.Error("Decode", zap.Error(err))
			http.Error(res, "Decode", http.StatusInternalServerError)
			return
		}
		if err := s.service.UpdateMetrics(ctx, batch); err != nil {
			s.logger.Error("UpdateMetrics", zap.Error(err))
			http.Error(res, "UpdateMetrics", http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write([]byte("Gauge successfully updated")); err != nil {
			s.logger.Error("res.Write", zap.Error(err))
			http.Error(res, "error res.Write", http.StatusInternalServerError)
			return
		}
	}
}

// UpdateMetricFromURL handles HTTP requests to update a metric from URL parameters.
func (s *restAPIServer) UpdateMetricFromURL(ctx context.Context) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var newMetrics model.Metrics
		if err := newMetrics.MapMetricsFromReqURL(req); err != nil {
			s.logger.Error("MapMetricsFromReq", zap.Error(err))
			http.Error(res, errInvalidRequestData.Error(), http.StatusBadRequest)
			return
		}
		if err := s.service.UpdateMetrics(ctx, []model.Metrics{newMetrics}); err != nil {
			s.logger.Error("UpdateMetrics", zap.Error(err))
			http.Error(res, "error res.Write", http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write([]byte("Gauge successfully updated")); err != nil {
			s.logger.Error("res.Write", zap.Error(err))
			http.Error(res, "error res.Write", http.StatusInternalServerError)
			return
		}
	}
}

// GetMetricFromJSON handles HTTP requests to retrieve a metric using JSON data.
func (s *restAPIServer) GetMetricFromJSON(ctx context.Context) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var newMetrics model.Metrics
		if err1 := json.NewDecoder(req.Body).Decode(&newMetrics); err1 != nil {
			s.logger.Error("Decode", zap.Error(err1))
			http.Error(res, err1.Error(), http.StatusBadRequest)
			return
		}
		metric, err2 := s.service.GetMetric(ctx, newMetrics)
		if err2 != nil {
			s.logger.Error("GetMetric", zap.Error(err2))
			http.Error(res, "GetMetric", http.StatusNotFound)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if err3 := json.NewEncoder(res).Encode(metric); err3 != nil {
			s.logger.Error("Encode", zap.Error(err3))
			http.Error(res, err3.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GetMetricFromURL handles HTTP requests to retrieve a metric using URL parameters.
func (s *restAPIServer) GetMetricFromURL(ctx context.Context) func(res http.ResponseWriter, _ *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var newMetrics model.Metrics
		if err := newMetrics.MapMetricsFromReqURL(req); err != nil {
			s.logger.Error("MapMetricsFromReq", zap.Error(err))
			http.Error(res, errInvalidRequestData.Error(), http.StatusBadRequest)
			return
		}
		metric, err2 := s.service.GetMetric(ctx, newMetrics)
		if err2 != nil {
			s.logger.Error("s.storage.GetMetricByModel", zap.Error(err2))
			http.Error(res, "error getting value from server", http.StatusNotFound)
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
			http.Error(res, "error res.Write", http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write([]byte(stringValue)); err != nil {
			s.logger.Error("res.Write", zap.Error(err))
			http.Error(res, "error res.Write", http.StatusInternalServerError)
			return
		}
	}
}

// GetMetrics handles HTTP requests to retrieve all metrics.
func (s *restAPIServer) GetMetrics(_ context.Context) func(res http.ResponseWriter, _ *http.Request) {
	return func(res http.ResponseWriter, _ *http.Request) {
		var ctx = context.Background()
		res.Header().Set("Content-Type", "text/html")
		body, err := s.service.GetAllMetrics(ctx)
		if err != nil {
			s.logger.Error("s.storage.GetAllMetrics", zap.Error(err))
			http.Error(res, "s.storage.GetAllMetrics", http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write([]byte(body)); err != nil {
			s.logger.Error("res.Write", zap.Error(err))
			http.Error(res, "error res.Write", http.StatusInternalServerError)
			return
		}
	}
}

// Ping handles HTTP requests to ping the server/database.
func (s *restAPIServer) Ping(ctx context.Context) func(res http.ResponseWriter, _ *http.Request) {
	return func(res http.ResponseWriter, _ *http.Request) {
		if s.config.DBEnable {
			newCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			if err := s.service.Ping(newCtx); err != nil {
				s.logger.Error("PingContext", zap.Error(err))
				http.Error(res, "data base is not responding", http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusOK)
			if _, err := res.Write([]byte("database is alive")); err != nil {
				s.logger.Error("res.Write", zap.Error(err))
				http.Error(res, "error res.Write", http.StatusInternalServerError)
				return
			}
		}
		res.WriteHeader(http.StatusInternalServerError)
	}
}
