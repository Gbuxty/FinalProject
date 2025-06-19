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
	Server   *ServerConfig
	Postgres *PostgresConfig
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
		Server:   newServerConfig(),
		Postgres: newPostgresConfig(),
	}
	return cfg, nil
}
