package server

import (
	"fmt"
	"gw-exchanger/internal/config"
	"gw-exchanger/internal/transport/handlers"
	"gw-exchanger/pkg/logger"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/Gbuxty/proto-exchange/exchange"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	server   *grpc.Server
	logger   logger.Logger
	handlers *handlers.ExchangerHandler
	cfg      *config.ServerConfig
}

const (
	serverConnTimeout = 10 * time.Second
)

func New(handlers *handlers.ExchangerHandler, cfg *config.ServerConfig, logger logger.Logger) *GrpcServer {
	return &GrpcServer{
		server:   grpc.NewServer(grpc.ConnectionTimeout(serverConnTimeout)),
		logger:   logger,
		handlers: handlers,
		cfg:      cfg,
	}
}

func (s *GrpcServer) Start() error {
	addr := fmt.Sprintf(":%d", s.cfg.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	pb.RegisterExchangeServiceServer(s.server, s.handlers)

	go func() {
		if err := s.server.Serve(lis); err != nil {
			s.logger.Errorf("failed to serve gRPC server: %v", err)
		}
	}()

	s.logger.Infof("gRPC server started port: %s", addr)
	return nil
}

func (s *GrpcServer) Stop() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.logger.Infof("Stopping gRPC server...")
	s.server.GracefulStop()

}
