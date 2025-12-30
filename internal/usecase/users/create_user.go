package user

import (
	"context"

	"main.go/internal/adapter/repository"
	"main.go/internal/domain/entity"
)

type CreateUserUseCase struct {
	repository repository.UserRepository
}

func NewCreateUserUseCase(repository repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repository: repository}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, payload *entity.User) error {
	return uc.repository.CreateUser(ctx, payload)
}
