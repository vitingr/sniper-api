package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"main.go/internal/domain/entity"
	"main.go/pkg/logger"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (repository *SessionRepository) CreateSession(ctx context.Context, payload *entity.Session) (*entity.Session, error) {
	query := `
		INSERT INTO sessions (id, user_email, refresh_token, is_revoked, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`

	session := &entity.Session{}
	err := repository.db.QueryRow(
		ctx, query,
		payload.ID, payload.UserEmail, payload.RefreshToken, payload.IsRevoked, payload.ExpiresAt,
	).Scan(&session.ID, &session.UserEmail, &session.RefreshToken, &session.IsRevoked, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		logger.Error("Cannot create a new session")
		return nil, err
	}

	return session, nil
}

func (repository *SessionRepository) GetSession(ctx context.Context, id string) (*entity.Session, error) {
	query := `SELECT * FROM sessions WHERE id = $1`

	session := &entity.Session{}
	err := repository.db.QueryRow(ctx, query, id).Scan(&session.ID, &session.UserEmail, &session.RefreshToken, &session.IsRevoked, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn("Session not found")
			return nil, errors.New("Session not found")
		}

		return nil, err
	}

	return session, nil
}

func (repository *SessionRepository) RevokeSession(ctx context.Context, id string) error {
	query := `UPDATE sessions SET is_revoked = TRUE WHERE id = $1`

	session := &entity.Session{}
	err := repository.db.QueryRow(ctx, query, id).Scan(&session.ID, &session.UserEmail, &session.RefreshToken, &session.IsRevoked, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn("Session not found")
			return errors.New("Session not found")
		}

		return err
	}

	return nil
}

func (repository *SessionRepository) DeleteSession(ctx context.Context, id string) error {
	query := `
		DELETE FROM sessions WHERE id = $1
		RETURNING *
	`

	session := &entity.Session{}
	err := repository.db.QueryRow(ctx, query, id).Scan(&session.ID, &session.UserEmail, &session.RefreshToken, &session.IsRevoked, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn("Session not found")
			return errors.New("Session not found")
		}

		return err
	}

	return nil
}
