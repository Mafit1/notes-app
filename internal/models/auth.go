package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID
	Email  string
	Role   string
	jwt.RegisteredClaims
}

type RefreshToken struct {
	TokenID   uuid.UUID `json:"token_id" db:"token_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	TokenHash string    `json:"-" db:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Revoked   bool      `json:"revoked" db:"revoked"`
	RevokedAt time.Time `json:"revoked_at" db:"revoked_at"`
}

func (rt *RefreshToken) IsExpired() bool {
	return rt.ExpiresAt.Before(time.Now())
}

func (rt *RefreshToken) IsActive() bool {
	return !rt.Revoked && !rt.IsExpired()
}

func (rt *RefreshToken) Revoke() {
	now := time.Now()
	rt.Revoked = true
	rt.RevokedAt = now
}
