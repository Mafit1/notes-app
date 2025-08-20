package jwtservice

import (
	"errors"
	"fmt"
	"time"

	"github.com/Mafit1/notes-app/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type service struct {
	secretKey []byte
	expiry    time.Duration
}

func New(secretKey string, expiry time.Duration) *service {
	if secretKey == "" {
		panic("JWT secret key is required")
	}
	if len(secretKey) < 32 {
		panic("JWT secret key must be at least 32 characters")
	}

	return &service{
		secretKey: []byte(secretKey),
		expiry:    expiry,
	}
}

func (s *service) GenerateToken(in RegisterIn) (tokenString string, err error) {
	claims := &models.Claims{
		UserID: in.UserID,
		Email:  in.Email,
		Role:   in.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "notes-app",
			Subject:   in.UserID.String(),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *service) ParseToken(tokenString string) (claims *models.Claims, err error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&models.Claims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.secretKey, nil
		},
		jwt.WithLeeway(5*time.Second),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("%w: malformed token", ErrInvalidToken)
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		if claims.UserID == uuid.Nil {
			return nil, fmt.Errorf("%w: midding user ID", ErrInvalidToken)
		}
		if claims.Email == "" {
			return nil, fmt.Errorf("%w: missing email", ErrInvalidToken)
		}
		return claims, nil
	}
	return nil, ErrInvalidToken
}

func (s *service) RefreshToken(oldTokenString string) (string, error) {
	claims, err := s.ParseToken(oldTokenString)
	if err != nil {
		return "", fmt.Errorf("failed to parse token for refresh: %w", err)
	}

	if time.Until(claims.ExpiresAt.Time) < time.Minute {
		return "", ErrTokenExpired
	}

	return s.GenerateToken(
		RegisterIn{
			UserID: claims.UserID,
			Email:  claims.Email,
			Role:   claims.Role,
		},
	)
}

func (s *service) ValidateToken(tokenString string) (bool, error) {
	_, err := s.ParseToken(tokenString)
	if err != nil {
		if errors.Is(err, ErrTokenExpired) || errors.Is(err, ErrTokenRevoked) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *service) GetRemainingTime(tokenString string) (time.Duration, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	return time.Until(claims.ExpiresAt.Time), nil
}
