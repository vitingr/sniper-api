package user

import (
	"context"

	"main.go/internal/adapter/repository"
	"main.go/internal/domain/entity"
)

type DeleteUserUseCase struct {
	repository repository.UserRepository
}

func NewDeleteUserUseCase(repository repository.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{repository: repository}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, email string) (*entity.User, error) {
	return uc.repository.DeleteUser(ctx, email)
}
