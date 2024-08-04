package grpc

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"

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

// RunServer starts server with the configured router.
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
