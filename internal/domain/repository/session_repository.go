package repository

import (
	"context"

	"main.go/internal/domain/entity"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, payload entity.Session) (*entity.Session, error)
	GetSession(ctx context.Context, id string) (*entity.Session, error)
	RevokeSession(ctx context.Context, id string) error
	DeleteSession(ctx context.Context, id string) error
}
