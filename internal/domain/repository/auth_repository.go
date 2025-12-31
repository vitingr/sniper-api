package repository

import (
	"context"

	"main.go/internal/adapter/repository"
	"main.go/internal/domain/entity"
	"main.go/pkg/jwt"
)

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	AccessToken string      `json:"access_token"`
	User        entity.User `json:"user"`
}

type AuthRepository interface {
	LoginUser(ctx context.Context, payload LoginUserRequest) (*LoginUserResponse, error)
	LogoutUser(ctx context.Context, session_id string) (error)
	SSORefreshToken(ctx context.Context, refresh_token string) (*repository.SSORefreshAccessTokenResponse, error)
	NewTokenMaker(ctx context.Context) (*jwt.Jwt, error)
	Me(ctx context.Context, access_token string) (*entity.User, error)
}
