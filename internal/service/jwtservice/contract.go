package jwtservice

import (
	"context"

	"github.com/Mafit1/notes-app/internal/models"
	"github.com/google/uuid"
)

type Service interface {
	GeneratePair(ctx context.Context, in GenerateIn) (*GenerateOut, error)
	RotatePair(ctx context.Context, oldTokenID uuid.UUID, in GenerateIn) (*GenerateOut, error)
	RefreshAccessToken(ctx context.Context, userID uuid.UUID, refreshToken string) (*GenerateOut, error)
	ValidateToken(tokenString string) (bool, error)
	ParseAccessToken(tokenString string) (*models.AccessTokenClaims, error)
}
