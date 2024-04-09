package service

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

func (s *Server) UpdateMetricFromJSON(ctx context.Context) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var newMetrics model.Metrics
		if err := newMetrics.MapMetricsFromReqJSON(req); err != nil {
			s.logger.Error("MapMetricsFromReqJSON", zap.Error(err))
			http.Error(res, errInvalidRequestData.Error(), http.StatusBadRequest)
			return
		}
		if err := s.storage.UpdateMetricValue(ctx, newMetrics); err != nil {
			s.logger.Error("UpdateMetricValue", zap.Error(err))
			http.Error(res, "error UpdateMetricValue", http.StatusInternalServerError)
		}
		if s.config.SyncStoreEnable {
			if err := s.storage.StoreMetrics(ctx, s.config.StoreFilePath); err != nil {
				s.logger.Error("StoreMetrics", zap.Error(err))
				http.Error(res, "error StoreMetrics", http.StatusInternalServerError)
			}
		}
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write([]byte("Gauge successfully updated")); err != nil {
			s.logger.Error("res.Write", zap.Error(err))
			http.Error(res, "error res.Write", http.StatusInternalServerError)
		}
	}
}

func (s *Server) UpdateMetricsFromJSON(ctx context.Context) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var batch []model.Metrics
		if err := json.NewDecoder(req.Body).Decode(&batch); err != nil {
			s.logger.Error("Decode", zap.Error(err))
			http.Error(res, "Decode", http.StatusInternalServerError)
		}

	}
}

func (s *Server) UpdateMetricFromURL(ctx context.Context) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var newMetrics model.Metrics
		if err := newMetrics.MapMetricsFromReqURL(req); err != nil {
			s.logger.Error("MapMetricsFromReq", zap.Error(err))
			http.Error(res, errInvalidRequestData.Error(), http.StatusBadRequest)
			return
		}
		if err := s.storage.UpdateMetricValue(ctx, newMetrics); err != nil {
			s.logger.Error("UpdateMetricValue", zap.Error(err))
			http.Error(res, "error UpdateMetricValue", http.StatusInternalServerError)
		}
		if s.config.SyncStoreEnable {
			if err := s.storage.StoreMetrics(ctx, s.config.StoreFilePath); err != nil {
				s.logger.Error("StoreMetrics", zap.Error(err))
				http.Error(res, "error StoreMetrics", http.StatusInternalServerError)
			}
		}
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write([]byte("Gauge successfully updated")); err != nil {
			s.logger.Error("res.Write", zap.Error(err))
			http.Error(res, "error res.Write", http.StatusInternalServerError)
		}
	}
}

func (s *Server) GetMetricFromJSON(ctx context.Context) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var newMetrics model.Metrics
		if err1 := json.NewDecoder(req.Body).Decode(&newMetrics); err1 != nil {
			s.logger.Error("Decode", zap.Error(err1))
			http.Error(res, err1.Error(), http.StatusBadRequest)
			return
		}
		metric, err2 := s.storage.GetMetricByModel(ctx, newMetrics)
		if err2 != nil {
			s.logger.Error("s.storage.GetMetricByModel", zap.Error(err2))
			http.Error(res, "error getting value from server", http.StatusNotFound)
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if err3 := json.NewEncoder(res).Encode(metric); err3 != nil {
			s.logger.Error("Encode", zap.Error(err3))
			http.Error(res, err3.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) GetMetricFromURL(ctx context.Context) func(res http.ResponseWriter, _ *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var newMetrics model.Metrics
		if err := newMetrics.MapMetricsFromReqURL(req); err != nil {
			s.logger.Error("MapMetricsFromReq", zap.Error(err))
			http.Error(res, errInvalidRequestData.Error(), http.StatusBadRequest)
			return
		}
		metric, err2 := s.storage.GetMetricByModel(ctx, newMetrics)
		if err2 != nil {
			s.logger.Error("s.storage.GetMetricByModel", zap.Error(err2))
			http.Error(res, "error getting value from server", http.StatusNotFound)
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
		}
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write([]byte(stringValue)); err != nil {
			s.logger.Error("res.Write", zap.Error(err))
			http.Error(res, "error res.Write", http.StatusInternalServerError)
		}
	}
}
func (s *Server) GetMetrics(_ context.Context) func(res http.ResponseWriter, _ *http.Request) {
	return func(res http.ResponseWriter, _ *http.Request) {
		var ctx = context.Background()
		res.Header().Set("Content-Type", "text/html")
		body, err := s.storage.GetAllMetrics(ctx)
		if err != nil {
			s.logger.Error("s.storage.GetAllMetrics", zap.Error(err))
			http.Error(res, "s.storage.GetAllMetrics", http.StatusInternalServerError)
		}
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write([]byte(body)); err != nil {
			s.logger.Error("res.Write", zap.Error(err))
			http.Error(res, "error res.Write", http.StatusInternalServerError)
		}
	}
}

func (s *Server) Ping(ctx context.Context) func(res http.ResponseWriter, _ *http.Request) {
	return func(res http.ResponseWriter, _ *http.Request) {
		if s.config.DBEnable {
			newCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			if err := s.db.PingContext(newCtx); err != nil {
				s.logger.Error("PingContext", zap.Error(err))
				http.Error(res, "data base is not responding", http.StatusInternalServerError)
			}
			res.WriteHeader(http.StatusOK)
			if _, err := res.Write([]byte("database is alive")); err != nil {
				s.logger.Error("res.Write", zap.Error(err))
				http.Error(res, "error res.Write", http.StatusInternalServerError)
			}
		}
		res.WriteHeader(http.StatusInternalServerError)
	}
}
