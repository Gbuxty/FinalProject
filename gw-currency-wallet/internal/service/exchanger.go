package service

import (
	"context"
	"fmt"
	"gw-currency-wallet/internal/adapter/kafka"
	"gw-currency-wallet/internal/adapter/redis"
	"gw-currency-wallet/internal/domain"
	"time"

	txmanager "gw-currency-wallet/internal/adapter/txManager"

	pb "github.com/Gbuxty/proto-exchange/exchange"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
)

const (
	amoutBigValue  = 30000.00
	eventTypeKafka = "wallet.transaction"
	redisTTL       = 5 * time.Hour
)

type ExchangerService struct {
	client     pb.ExchangeServiceClient
	repoWallet WalletOperation
	redis      *redis.Client
	tx         *txmanager.TxManager
	kafka      *kafka.Producer
}

func NewExchangerService(conn *grpc.ClientConn, repoWallet WalletOperation, redis *redis.Client, kafka *kafka.Producer, tx *txmanager.TxManager) *ExchangerService {
	return &ExchangerService{
		client:     pb.NewExchangeServiceClient(conn),
		repoWallet: repoWallet,
		redis:      redis,
		kafka:      kafka,
		tx:         tx,
	}
}

func (s *ExchangerService) GetAllRates(ctx context.Context) (map[string]float32, error) {

	resp, err := s.client.GetExchangeRates(ctx, &pb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rates: %w", err)
	}

	for key, value := range resp.Rates {
		if err := s.redis.SetRatesIntoCache(ctx, key, value, redisTTL); err != nil {
			return nil, fmt.Errorf("failed to cache rate %s: %w", key, err)
		}
	}

	return resp.Rates, nil
}

func (s *ExchangerService) ChangeCurrency(ctx context.Context, userID uuid.UUID, req domain.ExchangeRequest) ([]domain.Wallet, error) {
	if req.Amount.GreaterThan(decimal.NewFromFloat(amoutBigValue)) {

		message := &domain.WalletEvent{
			EventID:    uuid.New().String(),
			EventType:  eventTypeKafka,
			TimeCreate: time.Now().UTC(),
			Payload: &domain.TransactionPayload{
				UserID:       userID,
				FromCurrency: req.FromCurrency,
				ToCurrency:   req.ToCurrency,
				Amount:       req.Amount.String(),
			},
		}

		if err := s.kafka.SendMessage(userID.String(), message); err != nil {
			return nil, fmt.Errorf("failed send to kafka:%w", err)
		}
	}

	cacheKey := fmt.Sprintf("%s_%s", req.FromCurrency, req.ToCurrency)
	cachedRate, err := s.redis.GetRatesFromCache(ctx, cacheKey)
	if err != nil || cachedRate.IsZero() {
		resp, err := s.client.GetExchangeRateForCurrency(ctx, &pb.CurrencyRequest{
			FromCurrency: req.FromCurrency,
			ToCurrency:   req.ToCurrency,
		})
		if err != nil {
			return nil, fmt.Errorf("gRPC failed: %w", err)
		}
		cachedRate = decimal.NewFromFloat32(resp.Rate)
		if err := s.redis.SetRatesIntoCache(ctx, cacheKey, resp.Rate, 5*time.Hour); err != nil {
			return nil, fmt.Errorf("failed to cache rate %s: %w", cacheKey, err)
		}

	}

	balances := make([]domain.Wallet, 0)

	if err := s.tx.Do(ctx, func(ctx context.Context) error {

		currentBalance, err := s.repoWallet.GetBalanceForCurrency(ctx, userID, req.FromCurrency)
		if err != nil {
			return fmt.Errorf("failed get balance for currency:%w", err)
		}

		if currentBalance.LessThan(req.Amount) {
			return domain.ErrInsufficientFunds
		}
		rate, err := decimal.NewFromString(cachedRate.String())
		if err != nil {
			return fmt.Errorf("failed parse string to decimal:%w", err)
		}
		exchangedAmout := req.Amount.Mul(rate)

		withdrawReq := &domain.Wallet{
			Currency: req.FromCurrency,
			Amount:   req.Amount,
		}

		depositReq := &domain.Wallet{
			Currency: req.ToCurrency,
			Amount:   exchangedAmout,
		}
		if err := s.repoWallet.WithdrawBalance(ctx, userID, withdrawReq); err != nil {
			return fmt.Errorf("withdraw failed: %w", err)
		}

		if err := s.repoWallet.DepositBalance(ctx, userID, depositReq); err != nil {
			return fmt.Errorf("deposit failed: %w", err)
		}

		balances, err = s.repoWallet.GetCurrentBalance(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get balances: %w", err)
		}
		return nil

	}); err != nil {
		return nil, fmt.Errorf("failed tx:%w", err)
	}

	return balances, nil
}
