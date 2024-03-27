package service

import (
	"encoding/json"
	"errors"
	"github.com/mrkovshik/yametrics/internal/model"
	"net/http"

	"go.uber.org/zap"
)

var invalidRequestData = errors.New("invalid request data")

func (s *Server) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	var newMetrics model.Metrics
	if err := json.NewDecoder(req.Body).Decode(&newMetrics); err != nil {
		s.Logger.Error("Decode", zap.Error(err))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if !newMetrics.ValidateMetrics() {
		s.Logger.Error("ValidateMetrics", zap.Error(invalidRequestData))
		http.Error(res, invalidRequestData.Error(), http.StatusBadRequest)
		return
	}
	s.Storage.UpdateMetricValue(newMetrics)
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write([]byte("Gauge successfully updated")); err != nil {
		s.Logger.Error("res.Write", zap.Error(err))
		http.Error(res, "error res.Write", http.StatusInternalServerError)
	}

}

func (s *Server) GetMetric(res http.ResponseWriter, req *http.Request) {

	var newMetrics model.Metrics
	res.Header().Set("Content-Type", "application/json")
	if err1 := json.NewDecoder(req.Body).Decode(&newMetrics); err1 != nil {
		s.Logger.Error("Decode", zap.Error(err1))
		http.Error(res, err1.Error(), http.StatusBadRequest)
		return
	}
	//TODO: перенести проверки в validator
	//if metricType != metrics.MetricTypeCounter && metricType != metrics.MetricTypeGauge {
	//	http.Error(res, "invalid metric type", http.StatusBadRequest)
	//	return
	//}
	metric, err2 := s.Storage.GetMetricValue(newMetrics)
	if err2 != nil {
		s.Logger.Error("s.Storage.GetMetricValue", zap.Error(err2))
		http.Error(res, "error getting value from server", http.StatusNotFound)
	}
	res.WriteHeader(http.StatusOK)
	if err3 := json.NewEncoder(res).Encode(metric); err3 != nil {
		s.Logger.Error("Encode", zap.Error(err3))
		http.Error(res, err3.Error(), http.StatusInternalServerError)
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
