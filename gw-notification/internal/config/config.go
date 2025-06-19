package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	pathConfigFile = "config.env"
	dotenv         = "dotenv"
)

type Config struct {
	Server *HTTPServer
	Kafka  *KafkaConfig
	Mongo  *MongoConfig
}

func New() (*Config, error) {
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "docker" {
		viper.AutomaticEnv()
	} else {
		viper.SetConfigFile(pathConfigFile)
		viper.SetConfigType(dotenv)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config: %v", err)
		}
	}
	cfg := &Config{
		Server: newHTTPServer(),
		Kafka:  newKafkaConfig(),
		Mongo:  newMongoConfig(),
	}
	return cfg, nil
}
