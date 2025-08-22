package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mafit1/notes-app/internal/models"
	users_repo "github.com/Mafit1/notes-app/internal/repository/users"
	user_validator "github.com/Mafit1/notes-app/pkg/uservalidator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	usersRepository users_repo.Repository
	userValidator   user_validator.UserValidator
}

func New(repo users_repo.Repository, userValidator user_validator.UserValidator) Service {
	return &service{repo, userValidator}
}

func (s *service) Create(ctx context.Context, user CreateUser) (id uuid.UUID, err error) {
	if err = s.userValidator.ValidatePassword(user.Password); err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrPasswordValidation, err)
	}

	emailExist, err := s.usersRepository.EmailExists(ctx, user.Email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrFailedToCheckEmail, err)
	}
	if emailExist {
		return uuid.Nil, fmt.Errorf("%w: user with email %s already exists", ErrUserAlreadyExists, user.Email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrPasswordHashingFailed, err)
	}

	id, err = s.usersRepository.Create(
		ctx,
		users_repo.CreateUser{
			Name:     user.Name,
			Email:    user.Email,
			Password: string(hashedPassword),
		},
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrCannotCreateUser, err)
	}
	return id, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (user models.User, err error) {
	user, err = s.usersRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, users_repo.ErrUserNotFound) {
			return models.User{}, fmt.Errorf("%w: %v", ErrUserNotFound, err)
		}
		return models.User{}, fmt.Errorf("%w: %v", ErrCannotGetUser, err)
	}
	return user, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (user models.User, err error) {
	user, err = s.usersRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, users_repo.ErrUserNotFound) {
			return models.User{}, fmt.Errorf("%w: %v", ErrUserNotFound, err)
		}
		return models.User{}, fmt.Errorf("%w: %v", ErrCannotGetUser, err)
	}
	return user, nil
}
