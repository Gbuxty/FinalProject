package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	Amount   decimal.Decimal
	Currency string
}

type WalletEvent struct {
	EventID    string      `json:"event_id"`
	EventType  string      `json:"event_type"`
	TimeCreate time.Time   `json:"timestamp"`
	Payload    interface{} `json:"payload"`
}

type TransactionPayload struct {
	UserID       uuid.UUID `json:"user_id"`
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	Amount       string    `json:"amount"`
}