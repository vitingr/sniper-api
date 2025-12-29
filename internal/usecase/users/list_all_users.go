package user

import (
	"context"

	"main.go/internal/adapter/repository"
	"main.go/internal/domain/entity"
)

type ListAllUsersUseCase struct {
	repo repository.UserRepository
}

func NewListAllUsersUseCase(repo repository.UserRepository) *ListAllUsersUseCase {
	return &ListAllUsersUseCase{repo: repo}
}

func (uc *ListAllUsersUseCase) Execute(ctx context.Context, email string) ([]*entity.User, error) {
	return uc.repo.ListAllUsers(ctx, email)
}
