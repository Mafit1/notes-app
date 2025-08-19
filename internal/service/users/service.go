package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mafit1/notes-app/internal/models"
	users_repo "github.com/Mafit1/notes-app/internal/repository/users"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	usersRepository users_repo.Repository
}

func New(repo users_repo.Repository) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, user CreateUser) (id uuid.UUID, err error) {
	if len(user.Password) < 8 {
		return uuid.Nil, fmt.Errorf("%w: password should be at least 8 characters", ErrPasswordTooShort)
	}

	emailExist, err := s.usersRepository.EmailExists(ctx, user.Email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrFailedToCheckEmail, err)
	}
	if emailExist {
		return uuid.Nil, fmt.Errorf("%w: user with this email already exists", ErrUserAlreadyExists)
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
