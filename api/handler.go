package api

import (
	"context"
	"net/http"

	service "github.com/mrkovshik/yametrics/internal/service/server"
)

// UpdateMetricFromJSONHandler returns an HTTP handler function that calls s.UpdateMetricFromJSON.
func UpdateMetricFromJSONHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	ctx := context.Background()
	return s.UpdateMetricFromJSON(ctx)
}

// UpdateMetricFromURLHandler returns an HTTP handler function that calls s.UpdateMetricFromURL.
func UpdateMetricFromURLHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	ctx := context.Background()
	return s.UpdateMetricFromURL(ctx)
}

// GetMetricFromJSONHandler returns an HTTP handler function that calls s.GetMetricFromJSON.
func GetMetricFromJSONHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	ctx := context.Background()
	return s.GetMetricFromJSON(ctx)
}

// UpdateMetricsFromJSONHandler returns an HTTP handler function that calls s.UpdateMetricsFromJSON.
func UpdateMetricsFromJSONHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	ctx := context.Background()
	return s.UpdateMetricsFromJSON(ctx)
}

// GetMetricFromURLHandler returns an HTTP handler function that calls s.GetMetricFromURL.
func GetMetricFromURLHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	ctx := context.Background()
	return s.GetMetricFromURL(ctx)
}

// GetMetricsHandler returns an HTTP handler function that calls s.GetMetrics.
func GetMetricsHandler(s *service.Server) func(http.ResponseWriter, *http.Request) {
	ctx := context.Background()
	return s.GetMetrics(ctx)
}

// Ping returns an HTTP handler function that calls s.Ping.
func Ping(s *service.Server) func(http.ResponseWriter, *http.Request) {
	ctx := context.Background()
	return s.Ping(ctx)
}
