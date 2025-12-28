package entity

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	Password       string     `json:"password"`
	ProfilePicture string     `json:"profile_picture"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) IsValid() error {
	if u.Name == "" {
		return errors.New("First name is required")
	}

	if u.Email == "" && !strings.Contains(u.Email, "@") {
		return errors.New("Invalid email format")
	}

	return nil
}
