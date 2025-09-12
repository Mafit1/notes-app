package notes

import (
	"context"

	"github.com/Mafit1/notes-app/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=mocks/mock_$GOFILE -package=mocks . Service
type Service interface {
	Create(ctx context.Context, note CreateNote) (int64, error)
	CreateByUserID(ctx context.Context, userID uuid.UUID, note CreateNote) (int64, error)
	GetAll(ctx context.Context) ([]models.Note, error)
	GetByID(ctx context.Context, userID uuid.UUID, noteID int64) (*models.Note, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Note, error)
	Delete(ctx context.Context, userID uuid.UUID, noteID int64) error
	Update(ctx context.Context, userID uuid.UUID, note models.Note) (*models.Note, error)
}
