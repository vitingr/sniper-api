package auth

import (
	"context"

	"main.go/internal/adapter/repository"
	"main.go/internal/domain/entity"
)

type MeUseCase struct {
	repository repository.AuthRepository
}

func NewMeUseCase(repository repository.AuthRepository) *MeUseCase {
	return &MeUseCase{repository: repository}
}

func (uc *MeUseCase) Execute(ctx context.Context, access_token string) (*entity.User, error) {
	return uc.repository.Me(ctx, access_token)
}
