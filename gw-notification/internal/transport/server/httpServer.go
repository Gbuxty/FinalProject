package server

import (
	"context"
	"errors"
	"fmt"

	"gw-notification/internal/config"
	"gw-notification/pkg/logger"
	"net/http"
)

type HttpServer struct {
	server *http.Server
	cfg    *config.HTTPServer
	logger logger.Logger
}

func New(cfg *config.HTTPServer, logger logger.Logger) *HttpServer {
	return &HttpServer{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		cfg:    cfg,
		logger: logger,
	}
}

func (s *HttpServer) Start()  {
	s.logger.Infof("Server listening on port: %d", s.cfg.Port)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Errorf("failed start server :%v", err)
		}
	}()
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	s.logger.Info("Server stopped...")
	return nil
}
