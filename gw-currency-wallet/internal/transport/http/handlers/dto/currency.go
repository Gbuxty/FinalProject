package dto

type ChangeCurrency struct {
	From   string  `json:"from_currency" binding:"required,oneof=USD RUB EUR"`
	To     string  `json:"to_currency" binding:"required,oneof=USD RUB EUR"`
	Amount float32 `json:"amount" binding:"required,gt=0"`
}
