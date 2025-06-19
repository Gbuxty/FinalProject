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

	consumer, err := kafka.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.ConsumerGroup, logger)
	if err != nil {
		return fmt.Errorf("failed init kafka consumer:%w", err)
	}
	defer consumer.Close()

	wctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// вот эти 3 гоурутины не завершаются красиво ,нужно сделать что бы шатдаун был норм
	svcNotifi := service.NewNotificationService(consumer, logger, dbMongo)
	go svcNotifi.Start(wctx)

	svr := server.New(cfg.Server, logger)
	go svr.Start()

	go consumer.Start(wctx)
	ctxShutdownServer, cancel := context.WithTimeout(context.Background(), cfg.Server.ServerShutDownTimeout)
	defer cancel()
	svr.Shutdown(ctxShutdownServer)

	return nil
}
