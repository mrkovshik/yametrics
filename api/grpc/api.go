package grpc

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/mrkovshik/yametrics/internal/model"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	grpc2 "google.golang.org/grpc"

	"github.com/mrkovshik/yametrics/api"
	config "github.com/mrkovshik/yametrics/internal/config/server"
	pb "github.com/mrkovshik/yametrics/proto"
)

// Server represents the server configuration and dependencies.
type Server struct {
	server  *grpc2.Server
	service api.Service
	config  *config.ServerConfig
	logger  *zap.SugaredLogger
	pb.UnimplementedUsersServer
}

func NewServer(service api.Service, config *config.ServerConfig, logger *zap.SugaredLogger, server *grpc2.Server) *Server {
	return &Server{
		server:                   server,
		service:                  service,
		config:                   config,
		logger:                   logger,
		UnimplementedUsersServer: pb.UnimplementedUsersServer{},
	}
}

// RunServer starts the HTTP server with the configured router.
func (s *Server) RunServer(stop chan os.Signal) error {

	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		return err
	}

	pb.RegisterUsersServer(s.server, s)

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		if err := s.server.Serve(listen); err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		<-stop
		s.server.GracefulStop()
		return nil
	})

	if err := g.Wait(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}

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
