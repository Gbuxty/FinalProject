package config

import "github.com/spf13/viper"

type RedisConfig struct {
	Addr     string
	Password string
}

func newRedis() *RedisConfig {
	return &RedisConfig{
		Addr:     viper.GetString("REDIS_ADDR"),
		Password: viper.GetString("REDIS_PASSWORD"),
	}
}
