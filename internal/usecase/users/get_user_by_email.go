package user

import (
	"context"

	"main.go/internal/adapter/repository"
	"main.go/internal/domain/entity"
)

type GetUserByEmailUseCase struct {
	repository repository.UserRepository
}

func NewGetUserByEmailUseCase(repository repository.UserRepository) *GetUserByEmailUseCase {
	return &GetUserByEmailUseCase{repository: repository}
}

func (uc *GetUserByEmailUseCase) Execute(ctx context.Context, email string) (*entity.User, error) {
	return uc.repository.GetUserByEmail(ctx, email)
}
