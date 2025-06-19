package dto

import (
	"gw-currency-wallet/internal/domain"

	"github.com/shopspring/decimal"
)

type WalletOperation struct {
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required,oneof=USD RUB EUR"`
}

func MapToWalletOperation(d WalletOperation) *domain.Wallet {
	return &domain.Wallet{
		Amount:   decimal.NewFromFloat(d.Amount),
		Currency: d.Currency,
	}
}
