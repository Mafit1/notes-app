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
