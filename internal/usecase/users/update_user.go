package user

import (
	"context"

	"main.go/internal/adapter/repository"
	"main.go/internal/domain/entity"
)

type UpdateUserUseCase struct {
	repository repository.UserRepository
}

func NewUpdateUserUseCase(repository repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{repository: repository}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, payload *entity.User) error {
	return uc.repository.UpdateUser(ctx, payload)
}
