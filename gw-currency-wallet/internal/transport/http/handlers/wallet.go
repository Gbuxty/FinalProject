package handlers

import (
	"errors"
	"fmt"
	"gw-currency-wallet/internal/domain"
	"gw-currency-wallet/internal/service"
	"gw-currency-wallet/internal/transport/http/handlers/dto"
	"gw-currency-wallet/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletHandlers struct {
	service *service.WalletService
	logger  logger.Logger
}

func NewWalletHandlers(service *service.WalletService, logger logger.Logger) *WalletHandlers {
	return &WalletHandlers{
		service: service,
		logger:  logger,
	}
}

func (h *WalletHandlers) GetBalance(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Errorf("failed failed parse userid:%v", err)
		c.JSON(http.StatusInternalServerError, Error{Message: "failed failed parse userid"})
		return
	}

	currentBalance, err := h.service.GetWalletBalance(c.Request.Context(), userID)
	if err != nil {
		h.logger.Errorf("failed get user balance:%v", err)
		c.JSON(http.StatusInternalServerError, Error{Message: "failed get balance"})
		return
	}
	c.JSON(http.StatusOK, currentBalance)
}

func (h *WalletHandlers) Withdraw(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Errorf("failed parse userid:%v", err)
		c.JSON(http.StatusInternalServerError, Error{Message: "failed parse userid"})
		return
	}
	var req dto.WalletOperation

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, Error{Message: "invalid amount or currency"})
		return
	}

	wallet := dto.MapToWalletOperation(req)
	balance, err := h.service.WithdrawOperation(c.Request.Context(), userID, wallet)
	if err != nil {
		if errors.Is(err, domain.ErrInsufficientFundsOrNotWalletFound) {
			h.logger.Errorf("Failed withdraw balance: %v", err)
			c.JSON(http.StatusBadRequest, Error{Message: "Insufficient funds or invalid amount"})
			return
		}
		h.logger.Errorf("internal server error: %v", err)
		c.JSON(http.StatusInternalServerError, Error{Message: "internal server error"})
		return
	}

	walletResponse := WalletResponse{
		Message: "Withdrawal successful",
		Balance: balance,
	}

	c.JSON(http.StatusAccepted, walletResponse)
}

func (h *WalletHandlers) Deposit(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Errorf("failed failed parse userid:%v", err)
		c.JSON(http.StatusInternalServerError, Error{Message: "failed failed parse userid"})
		return
	}
	
	var req dto.WalletOperation

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, Error{Message: "invalid amount or currency"})
		return
	}
	wallet := dto.MapToWalletOperation(req)
	balance, err := h.service.DepositOperation(c.Request.Context(), userID, wallet)
	if err != nil {
		h.logger.Errorf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, Error{Message: "invalid amount or currency"})
		return
	}
	walletResponse := &WalletResponse{
		Message: "Account topped up successfully",
		Balance: balance,
	}
	c.JSON(http.StatusAccepted, walletResponse)
}

func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("user not authorized")
	}

	userIDstr, ok := userID.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("failed to parse user ID")
	}

	userIDuuid, err := uuid.Parse(userIDstr)
	if err != nil {
		return uuid.Nil, err
	}

	return userIDuuid, nil
}
