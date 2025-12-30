package session

import (
	"context"

	"main.go/internal/adapter/repository"
)

type DeleteSessionUseCase struct {
	repository repository.SessionRepository
}

func NewDeleteSessionUseCase(repository repository.SessionRepository) *DeleteSessionUseCase {
	return &DeleteSessionUseCase{repository: repository}
}

func (uc *DeleteSessionUseCase) Execute(ctx context.Context, id string) error {
	return uc.repository.DeleteSession(ctx, id)
}
