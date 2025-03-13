package auth

import (
	"context"
	"go-gin-e-comm/common"
	"go-gin-e-comm/user"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	userRepo  user.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo user.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (as *AuthService) Register(c context.Context, email, password string) error {
	userFound, err := as.userRepo.FindByEmail(c, email)
	if userFound != nil {
		return ErrUserExists
	}
	if err != nil {
		return common.ErrDatabase
	}
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	if err := as.userRepo.CreateUser(c, email, hashedPassword); err != nil {
		return ErrCreateUser
	}
	return nil
}

func (as *AuthService) Login(c context.Context, email, password string) (string, error) {
	userDB, _ := as.userRepo.FindByEmail(c, email)

	if userDB != nil && PasswordValid(userDB.Password, password) {
		return as.generateToken(userDB.ID)
	} else if userDB == nil || !PasswordValid(userDB.Password, password) {
		return "", ErrInvalidCredentials
	} else {
		return "", common.ErrDatabase
	}
}

func (as *AuthService) generateToken(userID primitive.ObjectID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(as.jwtSecret))
	if err != nil {
		slog.Error("Failed to generate token", slog.String("error", err.Error()))
		return "", ErrGenerateToken
	}

	return token, nil
}
