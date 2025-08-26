package auth

import "errors"

var (
	ErrRegistration          = errors.New("registration failed")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrPasswordValidation    = errors.New("password validation failed")
	ErrPasswordHashingFailed = errors.New("password hashing failed")
	ErrEmailValidation       = errors.New("email validation failed")
	ErrTokenGenerationFailed = errors.New("failed to generate token")
	ErrInvalidCredentials    = errors.New("invalid credentials")
)
