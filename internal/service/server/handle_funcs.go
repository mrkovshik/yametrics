package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mrkovshik/yametrics/internal/model"

	"go.uber.org/zap"
)

var errInvalidRequestData = errors.New("invalid request data")

func (s *Server) UpdateMetricFromJSON(res http.ResponseWriter, req *http.Request) {
	var newMetrics model.Metrics

	if err := newMetrics.MapMetricsFromReqJSON(req); err != nil {
		s.logger.Error("MapMetricsFromReqJSON", zap.Error(err))
		http.Error(res, errInvalidRequestData.Error(), http.StatusBadRequest)
		return
	}
	s.storage.UpdateMetricValue(newMetrics)
	if s.config.SyncStoreEnable {
		if err := s.storage.StoreMetrics(s.config.StoreFilePath); err != nil {
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

func (s *Server) UpdateMetricFromURL(res http.ResponseWriter, req *http.Request) {
	var newMetrics model.Metrics
	if err := newMetrics.MapMetricsFromReqURL(req); err != nil {
		s.logger.Error("MapMetricsFromReq", zap.Error(err))
		http.Error(res, errInvalidRequestData.Error(), http.StatusBadRequest)
		return
	}
	s.storage.UpdateMetricValue(newMetrics)
	if s.config.SyncStoreEnable {
		if err := s.storage.StoreMetrics(s.config.StoreFilePath); err != nil {
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

func (s *Server) GetMetricFromJSON(res http.ResponseWriter, req *http.Request) {
	var newMetrics model.Metrics

	if err1 := json.NewDecoder(req.Body).Decode(&newMetrics); err1 != nil {
		s.logger.Error("Decode", zap.Error(err1))
		http.Error(res, err1.Error(), http.StatusBadRequest)
		return
	}
	metric, err2 := s.storage.GetMetricByModel(newMetrics)
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

func (s *Server) GetMetricFromURL(res http.ResponseWriter, req *http.Request) {
	var newMetrics model.Metrics
	if err := newMetrics.MapMetricsFromReqURL(req); err != nil {
		s.logger.Error("MapMetricsFromReq", zap.Error(err))
		http.Error(res, errInvalidRequestData.Error(), http.StatusBadRequest)
		return
	}
	metric, err2 := s.storage.GetMetricByModel(newMetrics)
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

func (s *Server) GetMetrics(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	body, err := s.storage.GetAllMetrics()
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
