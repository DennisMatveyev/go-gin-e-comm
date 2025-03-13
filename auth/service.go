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
	log       *slog.Logger
}

func NewAuthService(userRepo user.UserRepository, jwtSecret string, log *slog.Logger) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret, log: log}
}

func (as *AuthService) Register(c context.Context, email, password string) error {
	userFound, err := as.userRepo.FindByEmail(c, email)
	if userFound != nil {
		return ErrUserExists
	}
	if err != nil {
		return common.ErrDatabase
	}
	hashedPassword, err := HashPassword(password, as.log)
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

	if userDB != nil && PasswordValid(userDB.Password, password, as.log) {
		return as.generateToken(userDB.ID)
	} else if userDB == nil || !PasswordValid(userDB.Password, password, as.log) {
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
		as.log.Error("Failed to generate token", "error", err.Error())
		return "", ErrGenerateToken
	}

	return token, nil
}
