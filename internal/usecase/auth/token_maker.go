package auth

import (
	"context"

	"main.go/internal/adapter/repository"
	"main.go/pkg/jwt"
)

type TokenMakerUseCase struct {
	repository repository.AuthRepository
}

func NewTokenMakerUseCase(repository repository.AuthRepository) *TokenMakerUseCase {
	return &TokenMakerUseCase{repository: repository}
}

func (uc *TokenMakerUseCase) Execute(ctx context.Context) (*jwt.Jwt) {
	return uc.repository.NewTokenMaker(ctx)
}
