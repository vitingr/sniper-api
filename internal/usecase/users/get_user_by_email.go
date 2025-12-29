package user

import (
	"context"

	"main.go/internal/adapter/repository"
	"main.go/internal/domain/entity"
)

type GetUserByEmailUseCase struct {
	repo repository.UserRepository
}

func NewGetUserByEmailUseCase(repo repository.UserRepository) *GetUserByEmailUseCase {
	return &GetUserByEmailUseCase{repo: repo}
}

func (uc *GetUserByEmailUseCase) Execute(ctx context.Context, email string) (*entity.User, error) {
	return uc.repo.GetUserByEmail(ctx, email)
}
