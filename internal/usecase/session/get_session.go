package session

import (
	"context"

	"main.go/internal/adapter/repository"
	"main.go/internal/domain/entity"
)

type GetSessionUseCase struct {
	repository repository.SessionRepository
}

func NewGetSessionUseCase(repository repository.SessionRepository) *GetSessionUseCase {
	return &GetSessionUseCase{repository: repository}
}

func (uc *GetSessionUseCase) Execute(ctx context.Context, id string) (*entity.Session, error) {
	return uc.repository.GetSession(ctx, id)
}
