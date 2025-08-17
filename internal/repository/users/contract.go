package users

import (
	"context"

	"github.com/Mafit1/notes-app/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
}
