package api

import (
	"net/http"

	service "github.com/mrkovshik/yametrics/internal/service/server"
)

func UpdateMetricFromJSONHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	return s.UpdateMetricFromJSON
}

func UpdateMetricFromURLHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	return s.UpdateMetricFromURL
}

func GetMetricFromJSONHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	return s.GetMetricFromJSON
}

func GetMetricFromURLHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	return s.GetMetricFromURL
}

func GetMetricsHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	return s.GetMetrics
}

func Ping(s *service.Server) func(http.ResponseWriter, *http.Request) {
	return s.Ping
}
