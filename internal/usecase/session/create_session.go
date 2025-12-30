package session

import (
	"context"

	"main.go/internal/adapter/repository"
	"main.go/internal/domain/entity"
)

type CreateSessionUseCase struct {
	repository repository.SessionRepository
}

func NewCreateSessionUseCase(repository repository.SessionRepository) *CreateSessionUseCase {
	return &CreateSessionUseCase{repository: repository}
}

func (uc *CreateSessionUseCase) Execute(ctx context.Context, payload *entity.Session) (*entity.Session, error) {
	return uc.repository.CreateSession(ctx, payload)
}
