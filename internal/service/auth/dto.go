package auth

import (
	"github.com/Mafit1/notes-app/internal/models"
	"github.com/google/uuid"
)

type RegisterIn struct {
	Name     string
	Email    string
	Password string
	Role     models.RoleType
}

type LoginIn struct {
	Email    string
	Password string
}

type RegisterOut struct {
	UserID       uuid.UUID
	Email        string
	AccessToken  string
	RefreshToken string
}

type LoginOut struct {
	UserID       uuid.UUID
	Email        string
	AccessToken  string
	RefreshToken string
}
