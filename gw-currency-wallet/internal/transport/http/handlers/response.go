package handlers

import "gw-currency-wallet/internal/domain"

type Error struct {
	Message string
}

type ResponseRegiser struct {
	Message string
}

type ResponseLogin struct {
	Message string
	Tokens  domain.TokensPair
}

type WalletResponse struct {
	Message string
	Balance []domain.Wallet
}
