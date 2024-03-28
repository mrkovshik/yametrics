package api

import (
	"net/http"

	service "github.com/mrkovshik/yametrics/internal/service/server"
)

func UpdateMetricHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	return s.UpdateMetric
}

func GetMetricHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	return s.GetMetric
}

func GetMetricsHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	return s.GetMetrics
}
