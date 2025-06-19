package service

import (
	"context"
	"fmt"

	"gw-exchanger/internal/domain"
	"gw-exchanger/pkg/logger"
)

type ExchangeService struct {
	repo   ExchangerRates
	logger logger.Logger
}

type ExchangerRates interface {
	GetRateForPair(ctx context.Context, baseCurrency, targetCurrency string) (*domain.ExchangeRate, error)
	GetAllRates(ctx context.Context) (map[string]float32, error)
}

func New(repo ExchangerRates, logger logger.Logger) *ExchangeService {
	return &ExchangeService{
		repo:   repo,
		logger: logger,
	}
}
func (s *ExchangeService) GetRates(ctx context.Context) (map[string]float32, error) {
	s.logger.Infof("Fetching all exchange rates")
	rates, err := s.repo.GetAllRates(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed get all rates:%w", err)
	}
	return rates, nil
}

func (s *ExchangeService) GetRateForOne(ctx context.Context, fromCurrency, toCurrency string) (*domain.ExchangeRate, error) {
	s.logger.Infof("Fetching rate: %s -> %s", fromCurrency, toCurrency)

	rates, err := s.repo.GetRateForPair(ctx, fromCurrency, toCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed get ratefor pair:%w", err)
	}
	return rates, nil
}
