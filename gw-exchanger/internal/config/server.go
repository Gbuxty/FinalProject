package config

import (
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port              int
	ConnectionTimeout time.Duration
}

func newServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:              viper.GetInt("SERVER_PORT"),
		ConnectionTimeout: viper.GetDuration("SERVER_CONNECTION_TIMEOUT"),
	}
}
