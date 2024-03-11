package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/mrkovshik/yametrics/internal/metrics"
	service "github.com/mrkovshik/yametrics/internal/service/server"
	"github.com/mrkovshik/yametrics/internal/storage/server"
	"net/http"
	"strconv"
)

func UpdateMetric(s *service.Server) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		metricName := chi.URLParam(req, "name")
		metricValue := chi.URLParam(req, "value")
		metricType := chi.URLParam(req, "type")
		switch metricType {

		case metrics.MetricTypeGauge:
			floatValue, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(res, "wrong value format", http.StatusBadRequest)
				return
			}
			gauge := server.NewGauge(metricName, floatValue)

			if err := gauge.Update(s.Storage); err != nil {
				http.Error(res, "Error updating counter", http.StatusBadRequest)
				return
			}

		case metrics.MetricTypeCounter:
			intValue, err := strconv.ParseInt(chi.URLParam(req, "value"), 0, 64)
			if err != nil {
				http.Error(res, "wrong value format", http.StatusBadRequest)
				return
			}
			counter := server.NewCounter(metricName, intValue)
			if err := counter.Update(s.Storage); err != nil {
				http.Error(res, "Error updating counter", http.StatusBadRequest)
				return
			}
		default:
			http.Error(res, "invalid metric type", http.StatusBadRequest)
			return
		}

		res.Write([]byte("Counter successfully updated"))
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
		switch metricType {

		case metrics.MetricTypeGauge:
			metricValue, err = s.Storage.GetGaugeValue(metricName)
			if err != nil {
				http.Error(res, "Data is missing", http.StatusNotFound)
			}

		case metrics.MetricTypeCounter:
			metricValue, err = s.Storage.GetCounterValue(metricName)
			if err != nil {
				http.Error(res, "Data is missing", http.StatusNotFound)
			}
		default:
			http.Error(res, "invalid metric type", http.StatusBadRequest)
			return
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
