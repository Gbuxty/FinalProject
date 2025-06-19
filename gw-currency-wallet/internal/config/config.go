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
	Server    *HTTPServer
	Postgres  *PostgresConfig
	Secret    *SecretKeyJWT
	Token     *TTlToken
	Redis     *RedisConfig
	Exchanger *ExchangerGrpc
	Kafka     *KafkaConfig
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
		Server:    newHTTPServer(),
		Postgres:  newPostgresConfig(),
		Secret:    newSecretKeyJWT(),
		Token:     newTTlToken(),
		Redis:     newRedis(),
		Exchanger: newExchangerServer(),
		Kafka:     newKafkaConfig(),
	}
	return cfg, nil
}
