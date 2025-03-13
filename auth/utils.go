package auth

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password", slog.String("error", err.Error()))
		return "", ErrHashPassword
	}
	return string(hashedBytes), nil
}

func PasswordValid(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		slog.Error("Failed to compare password", slog.String("error", err.Error()))
	}
	return err == nil
}
