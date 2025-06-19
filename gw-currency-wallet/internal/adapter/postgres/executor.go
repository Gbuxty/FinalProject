package postgres

import (
	"context"
	txmanager "gw-currency-wallet/internal/adapter/txManager"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Executor struct {
	pool *pgxpool.Pool
}

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewExecutor(pool *pgxpool.Pool) *Executor {
	return &Executor{pool: pool}
}

func (e *Executor) Get(ctx context.Context) DB {
	if tx, ok := ctx.Value(txmanager.TxKeyName).(pgx.Tx); ok {
		return tx 
	}
	return e.pool
}
