package jwtservice

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/Mafit1/notes-app/internal/models"
	auth_repo "github.com/Mafit1/notes-app/internal/repository/auth"
	users_service "github.com/Mafit1/notes-app/internal/service/users"
	"github.com/Mafit1/notes-app/pkg/hasher"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type service struct {
	accessSecret  []byte
	accessExpiry  time.Duration
	refreshSecret []byte
	refreshExpiry time.Duration
	authRepo      auth_repo.Repository
	usersService  users_service.Service
	hasher        hasher.Hasher
}

func New(
	accessSecret string,
	accessExpiry time.Duration,
	refreshSecret string,
	refreshExpiry time.Duration,
	authRepo auth_repo.Repository,
	usersService users_service.Service,
	hasher hasher.Hasher,
) Service {
	if len(accessSecret) < 32 || len(refreshSecret) < 32 {
		panic("JWT secret keys must be at least 32 characters")
	}

	return &service{
		accessSecret:  []byte(accessSecret),
		accessExpiry:  accessExpiry,
		refreshSecret: []byte(refreshSecret),
		refreshExpiry: refreshExpiry,
		authRepo:      authRepo,
		usersService:  usersService,
		hasher:        hasher,
	}
}

func (s *service) GeneratePair(ctx context.Context, in GenerateIn) (*GenerateOut, error) {
	accessToken, err := s.generateJWT(in, s.accessSecret, s.accessExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenPlain := uuid.New().String()
	h := hmac.New(sha256.New, s.refreshSecret)
	h.Write([]byte(refreshTokenPlain))
	refreshTokenHash := hex.EncodeToString(h.Sum(nil))

	_, err = s.authRepo.Create(ctx, auth_repo.RefreshTokenIn{
		UserID:    in.UserID,
		TokenHash: refreshTokenHash,
		TTL:       s.refreshExpiry,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	out := GenerateOut{
		AccessToken:      accessToken,
		RefreshToken:     refreshTokenPlain,
		RefreshExpiresAt: time.Now().Add(s.refreshExpiry),
	}

	return &out, nil
}

func (s *service) RotatePair(ctx context.Context, oldTokenID uuid.UUID, in GenerateIn) (*GenerateOut, error) {
	accessToken, err := s.generateJWT(in, s.accessSecret, s.accessExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshPlain := uuid.New().String()
	h := hmac.New(sha256.New, s.refreshSecret)
	h.Write([]byte(newRefreshPlain))
	newRefreshHash := hex.EncodeToString(h.Sum(nil))

	_, err = s.authRepo.Create(ctx, auth_repo.RefreshTokenIn{
		UserID:    in.UserID,
		TokenHash: newRefreshHash,
		TTL:       s.refreshExpiry,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save new refresh token: %w", err)
	}

	if err := s.authRepo.Revoke(ctx, oldTokenID); err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	return &GenerateOut{
		AccessToken:      accessToken,
		RefreshToken:     newRefreshPlain,
		RefreshExpiresAt: time.Now().Add(s.refreshExpiry),
	}, nil
}

func (s *service) RefreshAccessToken(ctx context.Context, refreshToken string) (*GenerateOut, error) {
	token, err := s.authRepo.GetByPlain(ctx, refreshToken, s.refreshSecret)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidRefreshToken, err)
	}

	if !token.IsActive() {
		return nil, fmt.Errorf("%w: token is not active", ErrInvalidRefreshToken)
	}

	user, err := s.usersService.GetByID(ctx, token.UserID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidRefreshToken, err)
	}

	return s.RotatePair(ctx, token.TokenID, GenerateIn{
		UserID: user.ID,
		Email:  user.Email,
		Role:   string(user.Role),
	})
}

func (s *service) RevokeRefreshToken(ctx context.Context, tokenPlain string) error {
	token, err := s.authRepo.GetByPlain(ctx, tokenPlain, s.refreshSecret)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidRefreshToken, err)
	}

	if !token.IsActive() {
		return fmt.Errorf("%w: token is already inactive", ErrTokenRevoked)
	}

	if err := s.authRepo.Revoke(ctx, token.TokenID); err != nil {
		return fmt.Errorf("failed to revoke token: %v", err)
	}

	return nil
}

func (s *service) ValidateToken(tokenString string) (bool, error) {
	_, err := s.parseJWT(tokenString, s.accessSecret)
	if err != nil {
		if errors.Is(err, ErrTokenExpired) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *service) ParseAccessToken(tokenString string) (*models.AccessTokenClaims, error) {
	return s.parseJWT(tokenString, s.accessSecret)
}

func (s *service) generateJWT(in GenerateIn, secret []byte, expiry time.Duration) (string, error) {
	claims := &models.AccessTokenClaims{
		UserID: in.UserID,
		Email:  in.Email,
		Role:   in.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "notes-app",
			Subject:   in.UserID.String(),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (s *service) parseJWT(tokenString string, secret []byte) (*models.AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&models.AccessTokenClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		},
		jwt.WithLeeway(5*time.Second),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if claims, ok := token.Claims.(*models.AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}
