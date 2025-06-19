package handlers

import (
	"errors"
	"gw-currency-wallet/internal/domain"

	"gw-currency-wallet/internal/service"
	"gw-currency-wallet/internal/transport/http/handlers/dto"

	"gw-currency-wallet/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type ExchangeHandlers struct {
	service *service.ExchangerService
	logger  logger.Logger
}

func NewExchangeHandlers(service *service.ExchangerService, logger logger.Logger) *ExchangeHandlers {
	return &ExchangeHandlers{service: service, logger: logger}
}

func (h *ExchangeHandlers) GetAllRates(c *gin.Context) {
	rates, err := h.service.GetAllRates(c.Request.Context())
	if err != nil {
		h.logger.Errorf("Failed to retrieve exchange rates: %v", err)
		c.JSON(http.StatusInternalServerError, Error{Message: "Failed to retrieve exchange rates"})
		return
	}
	h.logger.Infof("Actual rates from service exchange : %v", rates)
	c.JSON(http.StatusOK, rates)
}

func (h *ExchangeHandlers) ChangeCurrency(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Errorf("failed failed parse userid:%v", err)
		c.JSON(http.StatusInternalServerError, Error{Message: "failed failed parse userid"})
		return
	}

	var req dto.ChangeCurrency
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, Error{Message: "Invalid input"})
		return
	}

	domainReq := domain.ExchangeRequest{
		FromCurrency: req.From,
		ToCurrency:   req.To,
		Amount:       decimal.NewFromFloat32(req.Amount),
	}

	resp, err := h.service.ChangeCurrency(c.Request.Context(), userID, domainReq)
	if err != nil {
		if errors.Is(err, domain.ErrInsufficientFunds) {
			h.logger.Errorf("insufficient funds: %v", err)
			c.JSON(http.StatusBadRequest, Error{Message: "Insufficient funds or invalid currencies"})
			return
		}
		h.logger.Errorf("Exchange failed: %v", err)
		c.JSON(http.StatusInternalServerError, Error{Message: "internal server error"})
		return
	}

	h.logger.Infof("Change Currency success current balance: %v", resp)

	c.JSON(http.StatusOK, resp)
}
