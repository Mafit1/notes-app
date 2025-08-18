package users

import "errors"

var (
	ErrDatabase = errors.New("database error")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrCannotUpdateUser  = errors.New("cannot update user")
)
