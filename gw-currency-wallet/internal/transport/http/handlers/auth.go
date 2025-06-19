package handlers

import (
	"errors"
	"fmt"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/domain"
	"gw-currency-wallet/internal/service"
	"gw-currency-wallet/internal/transport/http/handlers/dto"
	"gw-currency-wallet/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthorizationService
	logger  logger.Logger
	cfg     *config.Config
}

const (
	Authorization = "Authorization"
)

func NewAuthHandler(service *service.AuthorizationService, logger logger.Logger, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
		cfg:     cfg,
	}
}
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterUser

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, Error{Message: "invalid name or password"})
		return
	}
	user := dto.MapToRegisterUser(req)
	if err := h.service.Register(c.Request.Context(), user); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			h.logger.Errorf("service user exists%s", err)
			c.JSON(http.StatusBadRequest, Error{Message: "user already exists"})
			return
		}
		h.logger.Errorf("service internal:%v", err)
		c.JSON(http.StatusInternalServerError, Error{Message: "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, ResponseRegiser{Message: "User registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginUser

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, Error{Message: "invalid name or password"})
		return
	}
	user := dto.MapToLoginUser(req)
	tokens, err := h.service.Login(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			h.logger.Errorf("user not found: %v", err)
			c.JSON(http.StatusUnauthorized, Error{Message: "invalid user or password"})
			return
		}
		h.logger.Errorf("service internal:%v", err)
		c.JSON(http.StatusInternalServerError, Error{Message: "internal server error"})
		return
	}
	bearer := fmt.Sprintf("Bearer %s", tokens.AccessToken)

	c.Header(Authorization, bearer)

	c.JSON(http.StatusOK, ResponseLogin{Message: "success login", Tokens: tokens})
}
