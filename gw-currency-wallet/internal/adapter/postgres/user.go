package postgres

import (
	"context"
	"errors"
	"fmt"
	"gw-currency-wallet/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	executor *Executor
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		executor: NewExecutor(pool),
	}
}

const (
	duplicatePostgresql = "23505"
)

func (r *UserRepository) SaveUser(ctx context.Context, req *domain.User) error {
	db := r.executor.Get(ctx)
	query := `INSERT INTO users (username, email, password_hash, created_at) VALUES ($1, $2, $3, NOW())`

	_, err := db.Exec(ctx, query, req.Name, req.Email, req.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == duplicatePostgresql {
			return domain.ErrUserAlreadyExists
		}
		return fmt.Errorf("failed save user error:%w", err)
	}
	return nil
}

func (r *UserRepository) GetUserAndPassordHashByName(ctx context.Context, req *domain.User) (domain.User, error) {
	db := r.executor.Get(ctx)
	query := `SELECT id,username,password_hash FROM users WHERE username=$1`
	var user domain.User
	err := db.QueryRow(ctx, query, req.Name).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("failed get username or passwordhash:%w", err)
	}
	return user, nil
}

func (r *UserRepository) SaveTokens(ctx context.Context, userID uuid.UUID, tokens *domain.TokensPair) error {
	db := r.executor.Get(ctx)
	query := `INSERT INTO tokens
	 						(
							 user_id,
	                         refresh_token,
	 						 access_token,
	 						 refresh_token_expires_at,
	  						 access_token_expires_at,
	 						 created_at) 
	  		  VALUES ($1, $2, $3,$4,$5, NOW())`

	_, err := db.Exec(
		ctx,
		query,
		userID,
		tokens.RefreshToken,
		tokens.AccessToken,
		tokens.RefreshTTL,
		tokens.AccessTTL,
	)
	if err != nil {
		return fmt.Errorf("failed save to db:%w", err)
	}

	return nil
}
