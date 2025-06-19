package postgres

import (
	"context"
	"errors"
	"fmt"
	"gw-currency-wallet/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type WalletRepository struct {
	pool *pgxpool.Pool
}

func NewWalletRepository(pool *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{pool: pool}
}

func (r *WalletRepository) DepositBalance(ctx context.Context, userID uuid.UUID, req *domain.Wallet) error {
	query := `INSERT INTO wallet(user_id,currency, balance, created_at, updated_at) 
			VALUES($1,$2,$3,NOW(),NOW())
			ON CONFLICT (user_id, currency) 
       		 DO UPDATE SET
		     balance = wallet.balance + EXCLUDED.balance,
            updated_at = NOW() `
	_, err := r.pool.Exec(ctx, query, userID, req.Currency, req.Amount)
	if err != nil {
		return fmt.Errorf("deposit balance: %w", err)
	}
	return nil
}

func (r *WalletRepository) WithdrawBalance(ctx context.Context, userID uuid.UUID, req *domain.Wallet) error {
	query := `
		UPDATE wallet 
		SET balance = balance - $1, updated_at = NOW() 
		WHERE user_id = $2 AND currency = $3 AND balance >= $1`
	res, err := r.pool.Exec(ctx, query, req.Amount, userID, req.Currency)
	if err != nil {
		return fmt.Errorf("withdraw funds: %w", err)
	}
	if res.RowsAffected() == 0 {
		return domain.ErrInsufficientFundsOrNotWalletFound
	}

	return nil
}

func (r *WalletRepository) GetCurrentBalance(ctx context.Context, userID uuid.UUID) ([]domain.Wallet, error) {
	query := `
		SELECT currency, balance 
		FROM wallet 
		WHERE user_id = $1  `

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query balances: %w", err)
	}
	defer rows.Close()

	var balances []domain.Wallet
	for rows.Next() {
		var w domain.Wallet
		if err := rows.Scan(&w.Currency, &w.Amount); err != nil {
			return nil, fmt.Errorf("scan balance row: %w", err)
		}
		balances = append(balances, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return balances, nil
}

func (r *WalletRepository) GetBalanceForCurrency(ctx context.Context, userID uuid.UUID, currency string) (decimal.Decimal, error) {
	var balanceStr string
	query := `SELECT balance FROM wallet WHERE user_id = $1 AND currency = $2`
	err := r.pool.QueryRow(ctx,
		query,
		userID,
		currency,
	).Scan(&balanceStr)

	if errors.Is(err, pgx.ErrNoRows) {
		return decimal.Zero, nil
	}

	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get balance: %w", err)
	}

	balance, err := decimal.NewFromString(balanceStr)
	if err != nil {
		return decimal.Zero, fmt.Errorf("invalid balance format: %w", err)
	}
	return balance, nil
}
