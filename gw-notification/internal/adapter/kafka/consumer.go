package kafka

import (
	"context"
	"fmt"
	"gw-notification/pkg/logger"

	"github.com/IBM/sarama"
)

type Consumer struct {
	group  sarama.ConsumerGroup
	topic  string
	logger logger.Logger
	msgCh  chan *sarama.ConsumerMessage
}

func NewKafkaConsumer(brokers []string, topic string, groupID string, logger logger.Logger) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest //??
	config.Version = sarama.V3_5_1_0

	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		logger.Errorf("error consumer group :%v", err)
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &Consumer{
		group:  group,
		topic:  topic,
		msgCh:  make(chan *sarama.ConsumerMessage),
		logger: logger,
	}, nil

}

func (c *Consumer) Start(ctx context.Context) {
	for {
		select {

		case <-ctx.Done():
			return
		default:
			if err := c.group.Consume(ctx, []string{c.topic}, c); err != nil {
				c.logger.Errorf("failed read from kafka comsume :%v", err)
			}

		}

	}
}

func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error {

	return nil
}

func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {

	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		c.msgCh <- message
		session.MarkMessage(message, "")
	}
	return nil
}

func (c *Consumer) MsgCh() <-chan *sarama.ConsumerMessage {
	return c.msgCh
}

func (c *Consumer) Close() error {
	if err := c.group.Close(); err != nil {
		return fmt.Errorf("failed close conn consumer group:%w", err)
	}
	return nil
}
