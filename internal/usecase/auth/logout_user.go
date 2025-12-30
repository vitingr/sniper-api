package auth

import (
	"context"

	"main.go/internal/adapter/repository"
)

type LogoutUserUseCase struct {
	repository repository.AuthRepository
}

func NewLogoutUserUseCase(repository repository.AuthRepository) *LogoutUserUseCase {
	return &LogoutUserUseCase{repository: repository}
}

func (uc *LogoutUserUseCase) Execute(ctx context.Context, session_id string) (error) {
	return uc.repository.LogoutUser(ctx, session_id)
}
