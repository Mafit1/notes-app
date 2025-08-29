package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mafit1/notes-app/internal/repository/auth"
	jwt_service "github.com/Mafit1/notes-app/internal/service/jwtservice"
	users_service "github.com/Mafit1/notes-app/internal/service/users"
	"github.com/Mafit1/notes-app/pkg/hasher"
	user_validator "github.com/Mafit1/notes-app/pkg/uservalidator"
	"github.com/google/uuid"
)

type service struct {
	authRepo      auth.Repository
	usersService  users_service.Service
	jwtService    jwt_service.Service
	userValidator user_validator.UserValidator
	hasher        hasher.Hasher
}

func New(authRepo auth.Repository, usersService users_service.Service, jwtService jwt_service.Service, userValidator user_validator.UserValidator, hasher hasher.Hasher) Service {
	return &service{
		authRepo:      authRepo,
		usersService:  usersService,
		jwtService:    jwtService,
		userValidator: userValidator,
		hasher:        hasher,
	}
}

func (s *service) Register(ctx context.Context, in RegisterIn) (out *RegisterOut, err error) {
	if err = s.userValidator.ValidateEmail(in.Email); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrEmailValidation, err)
	}

	if err = s.userValidator.ValidatePassword(in.Password); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPasswordValidation, err)
	}

	existing, err := s.usersService.GetByEmail(ctx, in.Email)
	if err != nil && !errors.Is(err, users_service.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := s.hasher.Hash(in.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPasswordHashingFailed, err)
	}

	userID, err := s.usersService.Create(
		ctx,
		users_service.CreateUser{
			Name:           in.Name,
			Email:          in.Email,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRegistration, err)
	}

	tokens, err := s.jwtService.GeneratePair(
		ctx,
		jwt_service.GenerateIn{
			UserID: userID,
			Email:  in.Email,
			Role:   string(in.Role),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenGenerationFailed, err)
	}

	return &RegisterOut{
		UserID:       userID,
		Email:        in.Email,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *service) Login(ctx context.Context, in LoginIn) (out *LoginOut, err error) {
	if err := s.userValidator.ValidateEmail(in.Email); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrEmailValidation, err)
	}

	if in.Password == "" {
		return nil, fmt.Errorf("%w: password cannot be empty", ErrPasswordValidation)
	}

	user, err := s.usersService.GetByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, users_service.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !s.hasher.Match(in.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	tokens, err := s.jwtService.GeneratePair(
		ctx,
		jwt_service.GenerateIn{
			UserID: user.ID,
			Email:  user.Email,
			Role:   string(user.Role),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenGenerationFailed, err)
	}

	return &LoginOut{
		UserID:       user.ID,
		Email:        user.Email,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *service) Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	token, err := s.authRepo.GetByPlain(ctx, refreshToken, userID, s.hasher)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidRefreshToken, err)
	}

	if !token.IsActive() {
		return fmt.Errorf("%w: token is already inactive", ErrCannotLogout)
	}

	if err := s.authRepo.Revoke(ctx, token.TokenID); err != nil {
		return fmt.Errorf("failed to revoke token: %v", err)
	}

	return nil
}

func (s *service) RevokeAll(ctx context.Context, userID uuid.UUID) (int64, error) {
	rowsAffected, err := s.authRepo.RevokeAllByUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to revoke all tokens: %w", err)
	}

	return rowsAffected, nil
}
