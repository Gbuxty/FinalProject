package postgresql

import (
	"context"
	"time"

	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	config.MinConns = 2
	config.MaxConns = 10
	config.MaxConnIdleTime = time.Hour
	config.HealthCheckPeriod = time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return pool, nil
}
