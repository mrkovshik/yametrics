package service

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrkovshik/yametrics/internal/metrics"
)

func (s *Server) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	metricName := chi.URLParam(req, "name")
	metricValue := chi.URLParam(req, "value")
	metricType := chi.URLParam(req, "type")
	if err := s.Storage.UpdateMetricValue(metricType, metricName, metricValue); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := res.Write([]byte("Gauge successfully updated")); err != nil {
		http.Error(res, "error res.Write", http.StatusInternalServerError)
	}
}

func (s *Server) GetMetric(res http.ResponseWriter, req *http.Request) {
	var (
		metricValue string
		err         error
	)
	metricName := chi.URLParam(req, "name")
	metricType := chi.URLParam(req, "type")
	if metricType != metrics.MetricTypeCounter && metricType != metrics.MetricTypeGauge {
		http.Error(res, "invalid metric type", http.StatusBadRequest)
		return
	}
	metricValue, err = s.Storage.GetMetricValue(metricType, metricName)
	if err != nil {
		http.Error(res, "error getting value from server", http.StatusNotFound)
	}
	if _, err := res.Write([]byte(metricValue)); err != nil {
		http.Error(res, "error res.Write", http.StatusInternalServerError)
	}
}

func (s *Server) GetMetrics(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	body, err := s.Storage.GetAllMetrics()
	if err != nil {
		http.Error(res, "s.Storage.GetAllMetrics", http.StatusInternalServerError)
	}
	if _, err := res.Write([]byte(body)); err != nil {
		http.Error(res, "error res.Write", http.StatusInternalServerError)
	}
}
