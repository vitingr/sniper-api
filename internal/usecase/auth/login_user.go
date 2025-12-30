package auth

import (
	"context"

	"main.go/internal/adapter/repository"
)

type LoginUserUseCase struct {
	repository repository.AuthRepository
}

func NewLoginUserUseCase(repository repository.AuthRepository) *LoginUserUseCase {
	return &LoginUserUseCase{repository: repository}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, payload *repository.LoginUserRequest) (*repository.LoginUserResponse, error) {
	return uc.repository.LoginUser(ctx, payload)
}
