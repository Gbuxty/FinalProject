package jwt

import (
	"fmt"
	"gw-currency-wallet/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SignedTokenAndExp struct {
	SignedToken string
	ExpiresAt   time.Time
}

func GenerateToken(req *domain.User, secretKey string, TokenTTL time.Duration) (SignedTokenAndExp, error) {
	expiresAt := time.Now().Add(TokenTTL)

	claims := jwt.MapClaims{
		"user_id": req.ID.String(),
		"email":   req.Email,
		"exp":     expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return SignedTokenAndExp{}, fmt.Errorf("failed to generate token: %w", err)
	}

	sigtoken := SignedTokenAndExp{
		SignedToken: signedToken,
		ExpiresAt:   expiresAt,
	}

	return sigtoken, nil
}

func CheckValidTokenAndGetUserID(tokenString, secretKey string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to parse token claims")
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token payload")
	}
	return userID, nil
}
