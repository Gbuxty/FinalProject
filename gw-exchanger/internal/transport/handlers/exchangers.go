package handlers

import (
	"context"
	"gw-exchanger/internal/service"
	"gw-exchanger/pkg/logger"

	pb "github.com/Gbuxty/proto-exchange/exchange"
)

type ExchangerHandler struct {
	service *service.ExchangeService
	logger  logger.Logger
	pb.UnimplementedExchangeServiceServer
}

func New(service *service.ExchangeService, logger logger.Logger) *ExchangerHandler {
	return &ExchangerHandler{
		service: service,
		logger:  logger,
	}
}


func (h *ExchangerHandler) GetExchangeRates(ctx context.Context, req *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	h.logger.Infof("Received request for exchange rates")

	rates, err := h.service.GetRates(ctx)
	if err != nil {
		h.logger.Errorf("Error fetching exchange rates: %v", err)
		return nil, err
	}

	response := &pb.ExchangeRatesResponse{
		Rates: rates,
	}
	
	return response, nil
}

func (h *ExchangerHandler) GetExchangeRateForCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	h.logger.Infof("Received request: %s -> %s", req.FromCurrency, req.ToCurrency)

	rate, err := h.service.GetRateForOne(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		h.logger.Errorf("Error fetching exchange rate: %v", err)
		return nil, err
	}

	return &pb.ExchangeRateResponse{
		FromCurrency: rate.FromCurrency,
		ToCurrency:   rate.ToCurrency,
		Rate:         rate.Rate,
	}, nil
}
