package config

import (
	"time"

	"github.com/spf13/viper"
)

type TTlToken struct {
	RefreshTTl time.Duration
	AccessTTl  time.Duration
}

func newTTlToken() *TTlToken {
	return &TTlToken{
		AccessTTl:  viper.GetDuration("ACCESS_TOKEN_TTL"),
		RefreshTTl: viper.GetDuration("REFRESH_TOKEN_TTL"),
	}
}
