package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"main.go/pkg/logger"
)

type UserClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewUserClaims(id string, email string, duration time.Duration) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()

	if err != nil {
		logger.Error("Error generating token ID")
		return nil, fmt.Errorf("Error generating token ID: %w", err)
	}

	return &UserClaims{
		Email: email,
		ID:    id,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}
