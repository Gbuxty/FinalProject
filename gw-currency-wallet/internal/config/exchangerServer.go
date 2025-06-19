package config

import "github.com/spf13/viper"

type ExchangerGrpc struct {
	Port string
}

func newExchangerServer() *ExchangerGrpc {
	return &ExchangerGrpc{
		Port: viper.GetString("EXCHANGER_SERVER_PORT"),
	}
}
