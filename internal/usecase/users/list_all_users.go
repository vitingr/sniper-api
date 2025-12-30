package user

import (
	"context"

	"main.go/internal/adapter/repository"
	"main.go/internal/domain/entity"
)

type ListAllUsersUseCase struct {
	repository repository.UserRepository
}

func NewListAllUsersUseCase(repository repository.UserRepository) *ListAllUsersUseCase {
	return &ListAllUsersUseCase{repository: repository}
}

func (uc *ListAllUsersUseCase) Execute(ctx context.Context, email string) ([]*entity.User, error) {
	return uc.repository.ListAllUsers(ctx, email)
}
