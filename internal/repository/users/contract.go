package users

import (
	"context"

	"github.com/Mafit1/notes-app/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=mocks/mock_$GOFILE -package=mocks . Repository
type Repository interface {
	Create(ctx context.Context, user CreateUser) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
}
