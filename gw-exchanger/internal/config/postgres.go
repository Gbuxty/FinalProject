package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func newPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     viper.GetString("POSTGRES_HOST"),
		Port:     viper.GetInt("POSTGRES_PORT"),
		User:     viper.GetString("POSTGRES_USER"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("POSTGRES_DB"),
		SSLMode:  viper.GetString("POSTGRES_SSL_MODE"),
	}
}
func (p *PostgresConfig) ToDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.DBName, p.SSLMode)
}
