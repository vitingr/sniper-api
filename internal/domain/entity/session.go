package entity

import (
	"errors"
	"time"
)

type Session struct {
	ID           string    `json:"id"`
	UserEmail    string    `json:"user_email"`
	RefreshToken string    `json:"refresh_token"`
	IsRevoked    bool      `json:"is_revoked"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (Session) TableName() string {
	return "sessions"
}

func (s *Session) IsValid() error {
	if s.UserEmail == "" {
		return errors.New("User email is required")
	}

	return nil
}
