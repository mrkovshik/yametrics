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
		s.Logger.Error("MapMetricsFromReqJSON", zap.Error(err))
		http.Error(res, errInvalidRequestData.Error(), http.StatusBadRequest)
		return
	}

	s.Storage.UpdateMetricValue(newMetrics)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write([]byte("Gauge successfully updated")); err != nil {
		s.Logger.Error("res.Write", zap.Error(err))
		http.Error(res, "error res.Write", http.StatusInternalServerError)
	}

}

func (s *Server) UpdateMetricFromURL(res http.ResponseWriter, req *http.Request) {
	var newMetrics model.Metrics
	if err := newMetrics.MapMetricsFromReqURL(req); err != nil {
		s.Logger.Error("MapMetricsFromReq", zap.Error(err))
		http.Error(res, errInvalidRequestData.Error(), http.StatusBadRequest)
		return
	}
	s.Storage.UpdateMetricValue(newMetrics)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write([]byte("Gauge successfully updated")); err != nil {
		s.Logger.Error("res.Write", zap.Error(err))
		http.Error(res, "error res.Write", http.StatusInternalServerError)
	}

}

func (s *Server) GetMetricFromJSON(res http.ResponseWriter, req *http.Request) {
	var newMetrics model.Metrics

	if err1 := json.NewDecoder(req.Body).Decode(&newMetrics); err1 != nil {
		s.Logger.Error("Decode", zap.Error(err1))
		http.Error(res, err1.Error(), http.StatusBadRequest)
		return
	}
	metric, err2 := s.Storage.GetMetricValue(newMetrics)
	if err2 != nil {
		s.Logger.Error("s.Storage.GetMetricValue", zap.Error(err2))
		http.Error(res, "error getting value from server", http.StatusNotFound)
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	if err3 := json.NewEncoder(res).Encode(metric); err3 != nil {
		s.Logger.Error("Encode", zap.Error(err3))
		http.Error(res, err3.Error(), http.StatusInternalServerError)
	}

}

func (s *Server) GetMetricFromURL(res http.ResponseWriter, req *http.Request) {
	var newMetrics model.Metrics
	if err := newMetrics.MapMetricsFromReqURL(req); err != nil {
		s.Logger.Error("MapMetricsFromReq", zap.Error(err))
		http.Error(res, errInvalidRequestData.Error(), http.StatusBadRequest)
		return
	}
	metric, err2 := s.Storage.GetMetricValue(newMetrics)
	if err2 != nil {
		s.Logger.Error("s.Storage.GetMetricValue", zap.Error(err2))
		http.Error(res, "error getting value from server", http.StatusNotFound)
	}

	var stringValue string
	switch metric.MType {
	case model.MetricTypeCounter:
		stringValue = fmt.Sprint(*metric.Delta)
	case model.MetricTypeGauge:
		stringValue = fmt.Sprint(*metric.Value)
	default:
		s.Logger.Error("invalid metric type", zap.Error(errors.New("ErrInvalidMetricType")))
		http.Error(res, "error res.Write", http.StatusInternalServerError)
	}
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write([]byte(stringValue)); err != nil {
		s.Logger.Error("res.Write", zap.Error(err))
		http.Error(res, "error res.Write", http.StatusInternalServerError)
	}
}

func (s *Server) GetMetrics(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	body, err := s.Storage.GetAllMetrics()
	if err != nil {
		s.Logger.Error("s.Storage.GetAllMetrics", zap.Error(err))
		http.Error(res, "s.Storage.GetAllMetrics", http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write([]byte(body)); err != nil {
		s.Logger.Error("res.Write", zap.Error(err))
		http.Error(res, "error res.Write", http.StatusInternalServerError)
	}
}
