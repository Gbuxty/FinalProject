package server

import (
	"context"
	"fmt"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type HttpServer struct {
	server *http.Server
	cfg    *config.HTTPServer
	logger logger.Logger
}

func New(cfg *config.HTTPServer, handler http.Handler, logger logger.Logger) *HttpServer {
	return &HttpServer{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		cfg:    cfg,
		logger: logger,
	}
}

func (s *HttpServer) Start() error {
	s.logger.Infof("Server listening on port: %d", s.cfg.Port)

	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	s.logger.Info("Server stopped...")
	return nil
}
