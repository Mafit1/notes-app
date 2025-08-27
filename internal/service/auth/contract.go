package auth

import (
	"context"

	"github.com/Mafit1/notes-app/internal/models"
)

type Service interface {
	Register(ctx context.Context, in RegisterIn) (*RegisterOut, error)
	Login(ctx context.Context, in LoginIn) (*LoginOut, error)
	Logout(ctx context.Context, token models.RefreshToken) error
}
