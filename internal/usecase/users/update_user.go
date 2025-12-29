package user

import (
	"context"

	"main.go/internal/adapter/repository"
	"main.go/internal/domain/entity"
)

type UpdateUserUseCase struct {
	repo repository.UserRepository
}

func NewUpdateUserUseCase(repo repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{repo: repo}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, input *entity.User) error {
	return uc.repo.UpdateUser(ctx, input)
}
