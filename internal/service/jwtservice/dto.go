package jwtservice

import (
	"time"

	"github.com/google/uuid"
)

type GenerateIn struct {
	UserID uuid.UUID
	Email  string
	Role   string
}

type GenerateOut struct {
	AccessToken      string
	RefreshToken     string
	RefreshExpiresAt time.Time
}
