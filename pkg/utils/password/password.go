package util

import (
	"golang.org/x/crypto/bcrypt"
	"main.go/pkg/logger"
)

func HashPassword(password string) (string, error) {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		logger.Error("Error hashing password")
		return "", err
	}

	return string(hashed_password), nil
}

func CheckPassword(password string, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
