package jwt

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"main.go/pkg/logger"
)

type Jwt struct {
	secret_key string
}

func NewJwt(secret_key string) *Jwt {
	return &Jwt{secret_key: secret_key}
}

func (maker *Jwt) CreateToken(id string, email string, duration time.Duration) (string, *UserClaims, error) {
	claims, err := NewUserClaims(id, email, duration)

	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(maker.secret_key))

	if err != nil {
		return "", nil, fmt.Errorf("Error signing token: %w", err)
	}

	return tokenString, claims, nil
}

func (maker *Jwt) VerifyToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			logger.Error("Invalid token signing method")
			return nil, fmt.Errorf("Invalid token signing method")
		}

		return []byte(maker.secret_key), nil
	})

	if err != nil {
		logger.Error("Error parsing JWT Token")
		return nil, fmt.Errorf("Error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)

	if !ok {
		logger.Error("Invalid token claims")
		return nil, fmt.Errorf("Invalid token claims")
	}

	return claims, nil
}
