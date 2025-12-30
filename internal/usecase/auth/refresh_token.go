package auth

import (
	"context"

	"main.go/internal/adapter/repository"
)

type RefreshTokenUseCase struct {
	repository repository.AuthRepository
}

func NewRefreshTokenUseCase(repository repository.AuthRepository) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{repository: repository}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, refresh_token string) (*repository.SSORefreshAccessTokenResponse, error) {
	return uc.repository.SSORefreshToken(ctx, refresh_token)
}
