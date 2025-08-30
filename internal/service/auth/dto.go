package auth

import (
	"time"

	"github.com/google/uuid"
)

type RegisterIn struct {
	Name     string
	Email    string
	Password string
}

type LoginIn struct {
	Email    string
	Password string
}

type AuthData struct {
	UserData    UserData
	AccessToken string
	RefreshData RefreshData
}

type UserData struct {
	ID    uuid.UUID
	Email string
}

type RefreshData struct {
	RefreshToken string
	ExpiresAt    time.Time
}
