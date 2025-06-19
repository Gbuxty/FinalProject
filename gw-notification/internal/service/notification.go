package service

import (
	"context"
	"encoding/json"
	"fmt"

	"gw-notification/internal/adapter/kafka"
	"gw-notification/internal/domain"
	"gw-notification/pkg/logger"

	"github.com/IBM/sarama"
)

type NotificationService struct {
	consumer *kafka.Consumer
	logger   logger.Logger
	repo     SaveEvent
	stop     chan struct{}
}

type SaveEvent interface {
	SaveEvent(ctx context.Context, req domain.Event) error
}

func NewNotificationService(consumer *kafka.Consumer, logger logger.Logger, repo SaveEvent) *NotificationService {
	return &NotificationService{
		consumer: consumer,
		logger:   logger,
		repo:     repo,
		stop:     make(chan struct{}),
	}
}

func (s *NotificationService) Start(ctx context.Context) {
	
	go func() {
		defer close(s.stop)
		for {
			select {
			case <-ctx.Done():
				s.logger.Infof("Notification service stopped")
				return
			case msg, ok := <-s.consumer.MsgCh():
				if !ok {
					return 
				}
				if err := s.processMessage(ctx, msg); err != nil {
					s.logger.Errorf("Failed to process message: %v", err)
				}
			}
		}
	}()

}

func (s *NotificationService) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var event domain.Event

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		s.logger.Errorf("Unmarshal json: %v", err)
		return err
	}
	s.logger.Infof("event: %v", event)
	if err := s.repo.SaveEvent(ctx, event); err != nil {
		return fmt.Errorf("failed save :%w", err)
	}

	s.logger.Infof("Message successfully saved:%v", event)
	return nil
}

func (s *NotificationService) Shutdown() {
	
	<-s.stop
	s.logger.Infof("Notification service shutdown complete")
}