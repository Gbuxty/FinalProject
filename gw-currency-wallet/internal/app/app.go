package app

import (
	"context"
	"fmt"
	"gw-currency-wallet/internal/adapter/kafka"
	"gw-currency-wallet/internal/adapter/postgres"
	"gw-currency-wallet/internal/adapter/redis"
	txmanager "gw-currency-wallet/internal/adapter/txManager"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/service"
	"gw-currency-wallet/internal/transport/http/handlers"
	"gw-currency-wallet/internal/transport/http/router"
	"gw-currency-wallet/internal/transport/server"
	grpcClient "gw-currency-wallet/pkg/client/grpc"
	"gw-currency-wallet/pkg/client/postgresql"
	"gw-currency-wallet/pkg/logger"
)

func Run() error {
	logger, err := logger.New()
	if err != nil {
		return fmt.Errorf("init logger:%s", err)
	}

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("init cfg:%s", err)
	}
	ctx := context.Background()

	db, err := postgresql.ConnectToDB(ctx, cfg.Postgres.ToDSN())
	if err != nil {
		logger.Errorf("failed to connect to database: %v", err)
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	txManager := txmanager.NewTxManager(db)

	redisClient := redis.New(cfg)
	defer redisClient.Close()

	addr := fmt.Sprintf(":%s", cfg.Exchanger.Port)
	conn, err := grpcClient.ConnectToExchanger(addr)
	if err != nil {
		logger.Errorf("failed connect to gRPC Exchanger: %v", err)
		return fmt.Errorf("failed connect to gRPC Exchanger: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Errorf("fail close connection:%v", err)
		}
	}()

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic, logger)
	if err != nil {
		return fmt.Errorf("failed init newProducer Kafka:%w", err)
	}
	defer producer.Close()

	repoUser := postgres.NewUserRepository(db)
	repoWallet := postgres.NewWalletRepository(db)

	svcUser := service.NewAuthorizationService(repoUser, cfg)
	svcWallet := service.NewWalletService(repoWallet)
	svcExchanger := service.NewExchangerService(conn, repoWallet, redisClient, producer, txManager)

	handlersUser := handlers.NewAuthHandler(svcUser, logger, cfg)
	walletHandler := handlers.NewWalletHandlers(svcWallet, logger)
	exchangerHandler := handlers.NewExchangeHandlers(svcExchanger, logger)

	router := router.New(logger, handlersUser, walletHandler, exchangerHandler)
	r := router.Routes()

	srv := server.New(cfg.Server, r, logger)
	go srv.Start()

	wctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ServerShutdownTimeout)
	defer cancel()
	srv.Shutdown(wctx)

	return nil
}
