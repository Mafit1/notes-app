package users

import "errors"

var (
	ErrCannotCreateUser  = errors.New("cannot create user")
	ErrCannotGetUser     = errors.New("cannot get user by id")
	ErrCannotUpdateUser  = errors.New("cannot update user")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
