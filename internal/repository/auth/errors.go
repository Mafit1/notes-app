package auth

import "errors"

var (
	ErrDatabase = errors.New("database error")

	ErrTokenNotFound       = errors.New("token not found")
	ErrCannotRevoke        = errors.New("cannot revoke token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
