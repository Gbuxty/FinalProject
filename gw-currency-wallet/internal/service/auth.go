package service

import (
	"context"
	"errors"
	"fmt"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/domain"
	"gw-currency-wallet/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthorizationService struct {
	repo UserAuthorization
	cfg  *config.Config
}

type UserAuthorization interface{
	SaveUser(ctx context.Context, req *domain.User) error
	GetUserAndPassordHashByName(ctx context.Context, req *domain.User) (domain.User, error)
	SaveTokens(ctx context.Context, userID uuid.UUID, tokens *domain.TokensPair) error
}

func NewAuthorizationService(repo UserAuthorization, cfg *config.Config) *AuthorizationService {
	return &AuthorizationService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *AuthorizationService) Register(ctx context.Context, req *domain.User) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generate password Hash:%w", err)
	}
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(passwordHash),
	}

	if err := s.repo.SaveUser(ctx, user); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return domain.ErrUserAlreadyExists
		}
		return fmt.Errorf("failed in service :%w", err)
	}

	return nil
}

func (s *AuthorizationService) Login(ctx context.Context, req *domain.User) (domain.TokensPair, error) {
	user, err := s.repo.GetUserAndPassordHashByName(ctx, req)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.TokensPair{}, domain.ErrUserNotFound
		}
		return domain.TokensPair{}, fmt.Errorf("failed service login:%w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return domain.TokensPair{}, fmt.Errorf("invalid password:%w", err)
	}

	access, err := jwt.GenerateToken(&user, s.cfg.Secret.Key, s.cfg.Token.AccessTTl)
	if err != nil {
		return domain.TokensPair{}, fmt.Errorf("failed generate jwt token:%w", err)
	}

	refresh, err := jwt.GenerateToken(&user, s.cfg.Secret.Key, s.cfg.Token.RefreshTTl)
	if err != nil {
		return domain.TokensPair{}, fmt.Errorf("failed generate jwt token:%w", err)
	}

	tokens := &domain.TokensPair{
		RefreshToken: domain.RefreshToken(refresh.SignedToken),
		AccessToken:  domain.AccessToken(access.SignedToken),
		RefreshTTL:   refresh.ExpiresAt,
		AccessTTL:    access.ExpiresAt,
	}

	if err := s.repo.SaveTokens(ctx, user.ID, tokens); err != nil {
		return domain.TokensPair{}, fmt.Errorf("failed save tokens:%w", err)
	}

	return *tokens, nil
}
