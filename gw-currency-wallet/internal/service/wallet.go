package service

import (
	"context"
	"errors"
	"fmt"
	"gw-currency-wallet/internal/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WalletService struct {
	repo WalletOperation
}

type WalletOperation interface {
	DepositBalance(ctx context.Context, userID uuid.UUID, req *domain.Wallet) error
	WithdrawBalance(ctx context.Context, userID uuid.UUID, req *domain.Wallet) error
	GetCurrentBalance(ctx context.Context, userID uuid.UUID) ([]domain.Wallet, error)
	GetBalanceForCurrency(ctx context.Context, userID uuid.UUID, currency string) (decimal.Decimal, error)
}

func NewWalletService(repo WalletOperation) *WalletService {
	return &WalletService{
		repo: repo,
	}
}

func (s *WalletService) GetWalletBalance(ctx context.Context, userID uuid.UUID) ([]domain.Wallet, error) {
	balance, err := s.repo.GetCurrentBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get balance: %w", err)
	}
	return balance, nil
}

func (s *WalletService) WithdrawOperation(ctx context.Context, userID uuid.UUID, req *domain.Wallet) ([]domain.Wallet, error) {
	err := s.repo.WithdrawBalance(ctx, userID, req)
	if err != nil {
		if errors.Is(err, domain.ErrInsufficientFundsOrNotWalletFound) {
			return nil, domain.ErrInsufficientFundsOrNotWalletFound
		}
		return nil, fmt.Errorf("failed to withdraw: %w", err)
	}

	balance, err := s.repo.GetCurrentBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get balance after withdraw: %w", err)
	}
	return balance, nil
}

func (s *WalletService) DepositOperation(ctx context.Context, userID uuid.UUID, req *domain.Wallet) ([]domain.Wallet, error) {
	err := s.repo.DepositBalance(ctx, userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to deposit: %w", err)
	}

	balance, err := s.repo.GetCurrentBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get balance after deposit: %w", err)
	}
	return balance, nil
}
