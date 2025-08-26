package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mafit1/notes-app/internal/models"
	users_repo "github.com/Mafit1/notes-app/internal/repository/users"
	"github.com/google/uuid"
)

type service struct {
	usersRepository users_repo.Repository
}

func New(repo users_repo.Repository) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, user CreateUser) (id uuid.UUID, err error) {
	emailExist, err := s.usersRepository.EmailExists(ctx, user.Email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrFailedToCheckEmail, err)
	}
	if emailExist {
		return uuid.Nil, fmt.Errorf("%w: user with email %s already exists", ErrUserAlreadyExists, user.Email)
	}

	id, err = s.usersRepository.Create(
		ctx,
		users_repo.CreateUser{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.HashedPassword,
		},
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrCannotCreateUser, err)
	}
	return id, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (user *models.User, err error) {
	user, err = s.usersRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, users_repo.ErrUserNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrCannotGetUser, err)
	}
	return user, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (user *models.User, err error) {
	user, err = s.usersRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, users_repo.ErrUserNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrCannotGetUser, err)
	}
	return user, nil
}
