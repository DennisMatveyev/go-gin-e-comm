package auth

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, log *slog.Logger) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password", "error", err.Error())
		return "", ErrHashPassword
	}
	return string(hashedBytes), nil
}

func PasswordValid(hashedPassword, password string, log *slog.Logger) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Error("Failed to compare password", "error", err.Error())
	}
	return err == nil
}
