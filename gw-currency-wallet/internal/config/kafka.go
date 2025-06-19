package config

import (
	"strings"

	"github.com/spf13/viper"
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

func newKafkaConfig() *KafkaConfig {
	brokers := strings.Split(viper.GetString("KAFKA_BROKERS"), ",")
	return &KafkaConfig{
		Brokers: brokers,
		Topic:   viper.GetString("KAFKA_TOPIC"),
	}
}
