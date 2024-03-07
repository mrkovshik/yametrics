package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/mrkovshik/yametrics/internal/metrics"
	"github.com/mrkovshik/yametrics/internal/service"
	"github.com/mrkovshik/yametrics/internal/storage"
	"net/http"
	"strconv"
)

func UpdateMetric(s *service.Service) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		metricName := chi.URLParam(req, "name")
		metricValue := chi.URLParam(req, "value")
		metricType := chi.URLParam(req, "type")
		switch metricType {

		case metrics.MetricTypeGauge:
			//if !verifyGaugeName(metricName) {
			//	http.Error(res, "Data is missing", http.StatusNotFound)
			//	return
			//}
			floatValue, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(res, "wrong value format", http.StatusBadRequest)
				return
			}
			gauge := storage.NewGauge(metricName, floatValue)

			if err := gauge.Update(s.Storage); err != nil {
				http.Error(res, "Error updating counter", http.StatusBadRequest)
				return
			}

		case metrics.MetricTypeCounter:
			//if metricName != "PollCount" && metricName != "testCounter" && metricName != "testSetGet197" {
			//	http.Error(res, "Data is missing", http.StatusNotFound)
			//	return
			//}
			intValue, err := strconv.ParseInt(chi.URLParam(req, "value"), 0, 64)
			if err != nil {
				http.Error(res, "wrong value format", http.StatusBadRequest)
				return
			}
			counter := storage.NewCounter(metricName, intValue)
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

func GetMetric(s *service.Service) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var metricValue string
		metricName := chi.URLParam(req, "name")
		metricType := chi.URLParam(req, "type")
		switch metricType {

		case metrics.MetricTypeGauge:
			//if !verifyGaugeName(metricName) {
			//	http.Error(res, "Data is missing", http.StatusNotFound)
			//	return
			//}
			metricValue = s.Storage.GetGaugeValue(metricName)

		case metrics.MetricTypeCounter:
			//if metricName != "PollCount" && metricName != "testCounter" {
			//	http.Error(res, "Data is missing", http.StatusNotFound)
			//	return
			//}
			metricValue = s.Storage.GetCounterValue(metricName)
		default:
			http.Error(res, "invalid metric type", http.StatusBadRequest)
			return
		}
		res.Write([]byte(metricValue))
	}
}

func GetMetrics(s *service.Service) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html")
		body := s.Storage.GetAllMetrics()
		res.Write([]byte(body))
	}
}

func verifyGaugeName(name string) bool {
	_, ok := metrics.MetricNamesMap[name]
	return ok
}
