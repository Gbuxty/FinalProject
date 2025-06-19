package config

import "github.com/spf13/viper"

type SecretKeyJWT struct {
	Key string
}

func newSecretKeyJWT() *SecretKeyJWT {
	return &SecretKeyJWT{
		Key: viper.GetString("SECRET_KEY_JWT"),
	}
}
