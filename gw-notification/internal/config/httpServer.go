package config

import (
	"time"

	"github.com/spf13/viper"
)

type HTTPServer struct {
	Port                  int
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	ServerShutDownTimeout time.Duration
}

func newHTTPServer() *HTTPServer {
	return &HTTPServer{
		Port:                  viper.GetInt("HTTP_SERVER_PORT"),
		ReadTimeout:           viper.GetDuration("HTTP_READ_TIMEOUT"),
		WriteTimeout:          viper.GetDuration("HTTP_WRITE_TIMEOUT"),
		ServerShutDownTimeout: viper.GetDuration("HTTP_SHUTDOWN_TIMEOUT"),
	}
}
