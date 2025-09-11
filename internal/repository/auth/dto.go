package auth

import (
	"time"

	"github.com/google/uuid"
)

type RefreshTokenIn struct {
	UserID    uuid.UUID
	TokenHash string
	TTL       time.Duration
}

type RefreshTokenOut struct {
	TokenID   uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	Revoked   bool
	RevokedAt *time.Time
}

func (rt *RefreshTokenOut) IsExpired() bool {
	return rt.ExpiresAt.Before(time.Now())
}

func (rt *RefreshTokenOut) IsActive() bool {
	return !rt.Revoked && !rt.IsExpired()
}

func (rt *RefreshTokenOut) Revoke() {
	now := time.Now()
	rt.Revoked = true
	rt.RevokedAt = &now
}
