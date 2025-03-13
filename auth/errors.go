package auth

import "errors"

var (
	ErrUserExists         = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrHashPassword       = errors.New("failed to hash password")
	ErrGenerateToken      = errors.New("could not generate token")
	ErrMissingAuthHeader  = errors.New("missing authorization header")
	ErrInvalidToken       = errors.New("invalid token")
	ErrCreateUser         = errors.New("failed to create user")
	ErrAdminRequired      = errors.New("admin access required")
	ErrUserNotFound       = errors.New("user not found")
)
