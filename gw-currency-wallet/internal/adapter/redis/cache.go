package redis

import (
	"context"
	"fmt"
	"gw-currency-wallet/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

type Client struct {
	client *redis.Client
}

func New(cfg *config.Config) *Client {
	return &Client{
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
		}),
	}
}

func (r *Client) GetRatesFromCache(ctx context.Context, key string) (decimal.Decimal, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return decimal.Zero, nil
	}
	if err != nil {
		return decimal.Zero, fmt.Errorf("redis get failed: %w", err)
	}
	return decimal.NewFromString(value)
}

func (r *Client) SetRatesIntoCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("redis set failed: %w", err)
	}
	return nil
}

func (r *Client) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis delete failed: %w", err)
	}

	return nil
}
func (r *Client) Close() error {
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("close rredis:%w", err)
	}
	return nil
}
