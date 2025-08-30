package auth

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, in RegisterIn) (*AuthData, error)
	Login(ctx context.Context, in LoginIn) (*AuthData, error)
	Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error
	RevokeAll(ctx context.Context, userID uuid.UUID) (int64, error)
}
