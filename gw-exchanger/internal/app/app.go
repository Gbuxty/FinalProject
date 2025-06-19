package app

import (
	"context"
	"fmt"
	"gw-exchanger/internal/adapter/postgres"
	"gw-exchanger/internal/config"
	"gw-exchanger/internal/server"
	"gw-exchanger/internal/service"
	"gw-exchanger/internal/transport/handlers"
	"gw-exchanger/pkg/client/postgresql"
	"gw-exchanger/pkg/logger"

	"log"
)

func Run() error {
	logger, err := logger.New()
	if err != nil {
		log.Fatal("logger not init")
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	cfg, err := config.New()
	if err != nil {
		logger.Errorf("config not init: %v", err)
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	ctx := context.Background()

	db, err := postgresql.ConnectToDB(ctx, cfg.Postgres.ToDSN())
	if err != nil {
		logger.Errorf("failed to connect to database: %v", err)
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()
	exchangeRepo := repository.New(db)
	exchangeService := service.New(exchangeRepo, logger)
	exchangeHandler := handlers.New(exchangeService, logger)

	srv := server.New(exchangeHandler, cfg.Server, logger)
	if err := srv.Start(); err != nil {
		logger.Errorf("failed to start server: %v", err)
		return fmt.Errorf("failed to start server: %w", err)
	}

	srv.Stop()
	return nil
}
