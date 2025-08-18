package notes

import (
	"context"

	"github.com/Mafit1/notes-app/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=mocks/mock_$GOFILE -package=mocks . Repository
type Repository interface {
	Create(ctx context.Context, note CreateNote) (int64, error)
	GetAll(ctx context.Context) ([]models.Note, error)
	GetAllFromUserByID(ctx context.Context, userID uuid.UUID) ([]models.Note, error)
	GetAllFromUserByEmail(ctx context.Context, userEmail string) ([]models.Note, error)
	GetByID(ctx context.Context, id int64) (models.Note, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, note models.Note) (updatedNote models.Note, err error)
}
