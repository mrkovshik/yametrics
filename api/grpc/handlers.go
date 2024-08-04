package grpc

import (
	"context"
	"errors"

	"github.com/mrkovshik/yametrics/internal/model"
	pb "github.com/mrkovshik/yametrics/proto"
)

func (s *Server) UpdateMetrics(ctx context.Context, request *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	mappedMetrics := make([]model.Metrics, len(request.Metrics))
	for i, metric := range request.Metrics {
		mappedMetrics[i] = model.Metrics{
			ID:    metric.ID,
			MType: metric.MType,
		}
		switch metric.MType {
		case model.MetricTypeGauge:
			mappedMetrics[i].Value = &metric.Value
		case model.MetricTypeCounter:
			mappedMetrics[i].Delta = &metric.Delta
		default:
			return &pb.UpdateMetricsResponse{Error: "unknown metric type"}, errors.New("unknown metric type")
		}

	}

	if err := s.service.UpdateMetrics(ctx, mappedMetrics); err != nil {
		return &pb.UpdateMetricsResponse{Error: err.Error()}, err
	}
	return &pb.UpdateMetricsResponse{}, nil
}
