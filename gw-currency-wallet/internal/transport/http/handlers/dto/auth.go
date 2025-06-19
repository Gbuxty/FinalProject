package dto

import "gw-currency-wallet/internal/domain"

type RegisterUser struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=20,excludesall= "`
	Password string `json:"password" binding:"required,min=8,max=32,containsany=!@#?*,excludesall= "`
	Email    string `json:"email" binding:"required,email"`
}
type LoginUser struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=20,excludesall= "`
	Password string `json:"password" binding:"required,min=8,max=32,containsany=!@#?*,excludesall= "`
}

func MapToRegisterUser(d RegisterUser) *domain.User {
	return &domain.User{
		Name:     d.Username,
		Email:    d.Email,
		Password: d.Password,
	}
}

func MapToLoginUser(d LoginUser) *domain.User {
	return &domain.User{
		Name:     d.Username,
		Password: d.Password,
	}
}
