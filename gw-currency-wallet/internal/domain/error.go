package domain

import "errors"

var (
	ErrUserNotFound                      = errors.New("user not found")
	ErrUserAlreadyExists                 = errors.New("user already exists")
	ErrInsufficientFundsOrNotWalletFound = errors.New("insufficient funds or wallet not found")
	ErrInsufficientFunds                 = errors.New("insufficient funds ")
)
