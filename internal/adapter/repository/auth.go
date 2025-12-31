package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"main.go/internal/domain/entity"
	"main.go/pkg/jwt"
	"main.go/pkg/logger"
	util "main.go/pkg/utils/password"
)

type AuthRepository struct {
	db         *pgxpool.Pool
	secret_key string
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	SessionID             string      `json:"session_id"`
	RefreshToken          string      `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time   `json:"refresh_token_expires_at"`
	AccessToken           string      `json:"access_token"`
	AccessTokenExpiresAt  time.Time   `json:"access_token_expires_at"`
	User                  entity.User `json:"user"`
}

type SSORefreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type SSORefreshAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func NewAuthRepository(db *pgxpool.Pool, secret_key string) *AuthRepository {
	return &AuthRepository{db: db, secret_key: secret_key}
}

func (repository *AuthRepository) LoginUser(ctx context.Context, payload *LoginUserRequest) (*LoginUserResponse, error) {
	query := `SELECT id, email, password FROM users WHERE email = $1`

	user := &entity.User{}
	err := repository.db.QueryRow(ctx, query, payload.Email).Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("User not found")
		}

		return nil, err
	}

	err = util.CheckPassword(payload.Password, user.Password)

	if err != nil {
		logger.Info("Wrong password at credentials login.")
		return &LoginUserResponse{}, err
	}

	token_maker := jwt.NewJwt(repository.secret_key)

	access_token, access_claims, err := token_maker.CreateToken(user.ID, user.Email, 15*time.Minute)
	if err != nil {
		logger.Error("Error creating JWT access token")
		return &LoginUserResponse{}, err
	}

	refresh_token, refresh_claims, err := token_maker.CreateToken(user.ID, user.Email, 7*24*time.Hour)
	if err != nil {
		logger.Error("Error creating JWT refresh token")
		return &LoginUserResponse{}, err
	}

	query = `
		INSERT INTO sessions (id, user_email, refresh_token, is_revoked, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`

	session := &entity.Session{}
	err = repository.db.QueryRow(
		ctx, query,
		refresh_claims.RegisteredClaims.ID, user.Email, refresh_token, false, refresh_claims.RegisteredClaims.ExpiresAt.Time,
	).Scan(&session.ID, &session.UserEmail, &session.RefreshToken, &session.IsRevoked, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		logger.Error("Cannot create a new session.")
		return nil, err
	}

	return &LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           access_token,
		RefreshToken:          refresh_token,
		AccessTokenExpiresAt:  access_claims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refresh_claims.ExpiresAt.Time,
		User: entity.User{
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (repository *AuthRepository) SSORefreshToken(ctx context.Context, refresh_token string) (*SSORefreshAccessTokenResponse, error) {
	token_maker := jwt.NewJwt(repository.secret_key)

	refresh_claims, err := token_maker.VerifyToken(refresh_token)
	if err != nil {
		logger.Error("Error varifying refresh token")
		return &SSORefreshAccessTokenResponse{}, err
	}

	query := `SELECT * FROM sessions WHERE id = $1`

	session := &entity.Session{}
	err = repository.db.QueryRow(ctx, query, refresh_claims.RegisteredClaims.ID).Scan(&session.ID, &session.UserEmail, &session.RefreshToken, &session.IsRevoked, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn("Session not found")
			return &SSORefreshAccessTokenResponse{}, errors.New("Session not found")
		}

		return &SSORefreshAccessTokenResponse{}, err
	}

	if session.IsRevoked {
		logger.Error("Session already revoked.")
		return &SSORefreshAccessTokenResponse{}, errors.New("Session already revoked.")
	}

	if session.UserEmail != refresh_claims.Email {
		logger.Error("Invalid session")
		return &SSORefreshAccessTokenResponse{}, errors.New("Invalid session")
	}

	access_token, access_claims, err := token_maker.CreateToken(refresh_claims.ID, refresh_claims.Email, 15*time.Minute)
	if err != nil {
		logger.Error("Error creating JWT access token")
		return &SSORefreshAccessTokenResponse{}, err
	}

	return &SSORefreshAccessTokenResponse{
		AccessToken:          access_token,
		AccessTokenExpiresAt: access_claims.ExpiresAt.Time,
	}, nil
}

func (repository *AuthRepository) LogoutUser(ctx context.Context, session_id string) error {
	query := `
		DELETE FROM sessions WHERE id = $1
		RETURNING *
	`

	session := &entity.Session{}
	err := repository.db.QueryRow(ctx, query, session_id).Scan(&session.ID, &session.UserEmail, &session.RefreshToken, &session.IsRevoked, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn("Session not found")
			return errors.New("Session not found")
		}

		return err
	}

	return nil
}

func (repository *AuthRepository) NewTokenMaker(ctx context.Context) *jwt.Jwt {
	return jwt.NewJwt(repository.secret_key)
}

func (repository *AuthRepository) Me(ctx context.Context, access_token string) (*entity.User, error) {
	token_maker := jwt.NewJwt(repository.secret_key)

	claims, err := token_maker.VerifyToken(access_token)
	if err != nil {
		logger.Error("Invalid access token.")
		return nil, err
	}

	query := `SELECT id, name, email, profile_picture FROM users WHERE id = $1`

	user := &entity.User{}
	err = repository.db.QueryRow(ctx, query, claims.RegisteredClaims.ID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.ProfilePicture,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("User not found.")
		}
		return nil, err
	}

	return user, nil
}
