package api

import (
	"context"
	"net/http"

	service "github.com/mrkovshik/yametrics/internal/service/server"
)

func UpdateMetricFromJSONHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	var ctx = context.Background()
	return s.UpdateMetricFromJSON(ctx)
}

func UpdateMetricFromURLHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	var ctx = context.Background()
	return s.UpdateMetricFromURL(ctx)
}

func GetMetricFromJSONHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	var ctx = context.Background()
	return s.GetMetricFromJSON(ctx)
}

func UpdateMetricsFromJSONHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	var ctx = context.Background()
	return s.UpdateMetricsFromJSON(ctx)
}

func GetMetricFromURLHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	var ctx = context.Background()
	return s.GetMetricFromURL(ctx)
}

func GetMetricsHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	var ctx = context.Background()
	return s.GetMetrics(ctx)
}

func Ping(s *service.Server) func(http.ResponseWriter, *http.Request) {
	ctx := context.Background()
	return s.Ping(ctx)
}
