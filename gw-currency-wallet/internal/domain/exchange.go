package domain

import "github.com/shopspring/decimal"

type ExchangeRequest struct {
	FromCurrency string
	ToCurrency   string
	Amount       decimal.Decimal
}

type ExchangeResponse struct {
	Message         string                    
	ExchangedAmount float32                   
	NewBalance      map[string]decimal.Decimal 
}
