package app

import (
	"context"
	"fmt"
	"gw-notification/internal/adapter/kafka"
	"gw-notification/internal/adapter/mongoDB"
	"gw-notification/internal/config"
	"gw-notification/internal/service"
	"gw-notification/internal/transport/server"
	"gw-notification/pkg/client/mongo"
	"gw-notification/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() error {
	logger, err := logger.New()
	if err != nil {
		return fmt.Errorf("failed init logger :%w", err)
	}

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed init config :%w", err)
	}
	ctxMongo, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.ConnectToMongo(ctxMongo, cfg.Mongo.URI, cfg.Mongo.DBName)
	if err != nil {
		return fmt.Errorf("failed connect mongo:%w", err)
	}
	defer client.Disconnect(ctxMongo)

	dbMongo := mongoDB.NewEventRepository(client.Database(cfg.Mongo.DBName), cfg.Mongo.Collection)

	wctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer, err := kafka.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.ConsumerGroup, logger)
	if err != nil {
		return fmt.Errorf("failed init kafka consumer:%w", err)
	}
	svcNotifi := service.NewNotificationService(consumer, logger, dbMongo)
	svr := server.New(cfg.Server, logger)

	svcNotifi.Start(wctx)
	consumer.Start(wctx)
	svr.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancel()

	svcNotifi.Shutdown()
	consumer.Shutdown()

	ctxShutdownServer, cancel := context.WithTimeout(context.Background(), cfg.Server.ServerShutDownTimeout)
	defer cancel()
	if err := svr.Shutdown(ctxShutdownServer); err != nil {
		return err
	}
	
	logger.Info("Shutdown complete")
	return nil
}
