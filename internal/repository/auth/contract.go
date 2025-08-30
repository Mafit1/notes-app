package auth

import (
	"context"

	"github.com/Mafit1/notes-app/pkg/hasher"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, in RefreshTokenIn) (uuid.UUID, error)
	GetByID(ctx context.Context, tokenID uuid.UUID) (*RefreshTokenOut, error)
	GetByHash(ctx context.Context, hash string) (*RefreshTokenOut, error)
	GetByPlain(ctx context.Context, plain string, userID uuid.UUID, hasher hasher.Hasher) (*RefreshTokenOut, error)
	GetAllByUser(ctx context.Context, userID uuid.UUID) ([]*RefreshTokenOut, error)
	Revoke(ctx context.Context, tokenID uuid.UUID) error
	RevokeAllByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	DeleteExpired(ctx context.Context) (int64, error)
}
