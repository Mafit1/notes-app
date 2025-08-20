package jwtservice

import (
	"time"

	"github.com/Mafit1/notes-app/internal/models"
)

type Service interface {
	GenerateToken(in RegisterIn) (string, error)
	ParseToken(tokenString string) (*models.Claims, error)
	RefreshToken(oldTokenString string) (string, error)
	ValidateToken(tokenString string) (bool, error)
	GetRemainingTime(tokenString string) (time.Duration, error)
}
