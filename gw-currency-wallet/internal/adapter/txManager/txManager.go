package txmanager

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TxKey string

const (
	TxKeyName TxKey = "pgx_tx"
)

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

func (t *TxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			if err = tx.Rollback(ctx); err != nil {
				err = fmt.Errorf("rollback tx: %w:%w", err, err)
				return
			}
			return
		}
		if e := tx.Commit(ctx); e != nil {
			err = fmt.Errorf("commit tx: %w:%w", err, e)
		}
	}()
	ctxWithTx := context.WithValue(ctx, TxKeyName, tx)
	if err := fn(ctxWithTx); err != nil {
		
		return err
	}
	return nil
}
