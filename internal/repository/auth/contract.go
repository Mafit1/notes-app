package auth

import (
	"context"

	"github.com/Mafit1/notes-app/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, in RefreshTokenIn) (uuid.UUID, error)
	GetByID(ctx context.Context, tokenID uuid.UUID) (*models.RefreshToken, error)
	GetByHash(ctx context.Context, hash string) (*models.RefreshToken, error)
	Revoke(ctx context.Context, tokenID uuid.UUID) error
	RevokeAllByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	DeleteExpired(ctx context.Context) (int64, error)
}
