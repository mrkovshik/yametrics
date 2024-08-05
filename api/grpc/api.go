package grpc

import (
	"context"
	"errors"
	"net"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	grpc2 "google.golang.org/grpc"

	"github.com/mrkovshik/yametrics/api"
	config "github.com/mrkovshik/yametrics/internal/config/server"
	pb "github.com/mrkovshik/yametrics/proto"
)

// Server represents a gRPC server.
type Server struct {
	server  *grpc2.Server
	service api.Service
	config  *config.ServerConfig
	logger  *zap.SugaredLogger
	pb.UnimplementedUsersServer
}

// NewServer creates a new Server instance.
//
// Parameters:
//   - service: The service instance implementing business logic.
//   - config: The server configuration.
//   - logger: The logger instance.
//   - server: The gRPC server instance.
//
// Returns:
//   - *Server: A new Server instance.
func NewServer(service api.Service, config *config.ServerConfig, logger *zap.SugaredLogger, server *grpc2.Server) *Server {
	return &Server{
		server:                   server,
		service:                  service,
		config:                   config,
		logger:                   logger,
		UnimplementedUsersServer: pb.UnimplementedUsersServer{},
	}
}

// RunServer starts the gRPC server and listens for shutdown signals.
//
// Parameters:
//   - stop: A channel to receive OS signals for graceful shutdown.
//
// Returns:
//   - error: An error if the server fails to start or stop gracefully.
func (s *Server) RunServer(ctx context.Context) error {
	// Listen on TCP port 3200
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		return err
	}

	// Register the Users service with the gRPC server
	pb.RegisterUsersServer(s.server, s)

	// Create an errgroup with background context
	g, _ := errgroup.WithContext(context.Background())

	// Start the gRPC server
	g.Go(func() error {
		if err := s.server.Serve(listen); err != nil {
			return err
		}
		return nil
	})

	// Wait for stop signal and gracefully stop the server
	g.Go(func() error {
		<-ctx.Done()
		s.server.GracefulStop()
		return nil
	})

	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
