package repository

import (
	"context"
	"fmt"
	"gw-exchanger/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExchangeRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *ExchangeRepository {
	return &ExchangeRepository{
		pool: pool,
	}
}

func (r *ExchangeRepository) GetRateForPair(ctx context.Context, baseCurrency, targetCurrency string) (*domain.ExchangeRate, error) {
	query := `SELECT rate FROM exchange_rates WHERE base_currency = $1 AND target_currency = $2`

	var rate float32
	err := r.pool.QueryRow(ctx, query, baseCurrency, targetCurrency).Scan(&rate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rate for %s -> %s: %w", baseCurrency, targetCurrency, err)
	}

	return &domain.ExchangeRate{
		FromCurrency: baseCurrency,
		ToCurrency:   targetCurrency,
		Rate:         rate,
	}, nil
}

func (r *ExchangeRepository) GetAllRates(ctx context.Context) (map[string]float32, error) {
	query := `SELECT base_currency, target_currency, rate FROM exchange_rates`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	rates := make(map[string]float32)

	for rows.Next() {
		var base, target string
		var rate float32

		if err := rows.Scan(&base, &target, &rate); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		key := fmt.Sprintf("%s_%s", base, target)
		rates[key] = rate

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}

	return rates, nil
}
