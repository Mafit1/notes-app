package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mafit1/notes-app/internal/models"
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

	refreshTokenPlain := uuid.New().String()
	refreshTokenHash, err := s.hasher.Hash(refreshTokenPlain)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenHashingFailed, err)
	}

	_, err = s.authRepo.Create(ctx, auth.RefreshTokenIn{
		UserID:    userID,
		TokenHash: refreshTokenHash,
		TTL: ,
	})

	AccessToken, err := s.jwtService.GenerateToken(jwt_service.GenerateIn{
		UserID: userID,
		Email:  in.Email,
		Role:   string(in.Role),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrAccessTokenGenerationFailed, err)
	}

	out = &RegisterOut{
		UserID:      userID,
		Email:       in.Email,
		AccessToken: AccessToken,
	}

	return out, nil
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

	authToken, err := s.jwtService.GenerateToken(jwt_service.GenerateIn{
		UserID: user.ID,
		Email:  user.Email,
		Role:   string(user.Role),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenGenerationFailed, err)
	}

	out = &LoginOut{
		UserID:    user.ID,
		Email:     user.Email,
		AuthToken: authToken,
	}

	return out, nil
}

func (s *service) Logout(ctx context.Context, token models.RefreshToken) error {
	panic("unimplemented")
}
