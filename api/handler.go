package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrkovshik/yametrics/internal/metrics"
	service "github.com/mrkovshik/yametrics/internal/service/server"
)

func UpdateMetric(s *service.Server) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		metricName := chi.URLParam(req, "name")
		metricValue := chi.URLParam(req, "value")
		metricType := chi.URLParam(req, "type")
		if err := s.Storage.UpdateMetricValue(metricType, metricName, metricValue); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		res.Write([]byte("Gauge successfully updated"))
	}

}

func GetMetric(s *service.Server) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
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
		res.Write([]byte(metricValue))
	}
}

func GetMetrics(s *service.Server) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html")
		body := s.Storage.GetAllMetrics()
		res.Write([]byte(body))
	}
}
