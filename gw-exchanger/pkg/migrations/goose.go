package migrations

import (
	"fmt"
	"gw-exchanger/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	migrationsDir = "./migrations"
	dialect       = "postgres"
)

func Run(pool *pgxpool.Pool, logger logger.Logger) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()
	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	logger.Info("Migrations applied successfully")
	return nil
}
