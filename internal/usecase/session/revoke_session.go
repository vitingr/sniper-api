package session

import (
	"context"

	"main.go/internal/adapter/repository"
)

type RevokeSessionUseCase struct {
	repository repository.SessionRepository
}

func NewRevokeSessionUseCase (repository repository.SessionRepository) *RevokeSessionUseCase {
	return &RevokeSessionUseCase{repository: repository}
}

func (uc *RevokeSessionUseCase) Execute(ctx context.Context, id string) error {
	return uc.repository.RevokeSession(ctx, id)
}
