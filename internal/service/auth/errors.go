package auth

import "errors"

var (
	ErrRegistration                = errors.New("registration failed")
	ErrUserAlreadyExists           = errors.New("user already exists")
	ErrUserNotFound                = errors.New("user not found")
	ErrPasswordValidation          = errors.New("password validation failed")
	ErrPasswordHashingFailed       = errors.New("password hashing failed")
	ErrTokenHashingFailed          = errors.New("token hashing failed")
	ErrEmailValidation             = errors.New("email validation failed")
	ErrAccessTokenGenerationFailed = errors.New("failed to generate access token")
	ErrCannotCreateRefreshToken    = errors.New("cannot create refresh token")
	ErrInvalidCredentials          = errors.New("invalid credentials")
)
