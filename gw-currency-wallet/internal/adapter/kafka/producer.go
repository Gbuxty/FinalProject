package kafka

import (
	"encoding/json"
	"fmt"
	"gw-currency-wallet/pkg/logger"

	"github.com/IBM/sarama"
)

type Producer struct {
	syncProducer sarama.SyncProducer
	topic        string
	logger       logger.Logger
}

func NewProducer(brokers []string, topic string, logger logger.Logger) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	prod, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("fail new producer:%w", err)
	}

	return &Producer{
		syncProducer: prod,
		topic:        topic,
		logger:       logger,
	}, nil
}

func (p *Producer) SendMessage(key string, message interface{}) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed marshal message:%w", err)
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(bytes),
	}

	partition, offset, err := p.syncProducer.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("unsuccessful send message producer kafka partition:%d,offset:%d", partition, offset)
	}
	p.logger.Infof("Success send message kafka partition,offset:", partition, offset)
	return nil
}

func (p *Producer) Close() error {
	return p.syncProducer.Close()
}
