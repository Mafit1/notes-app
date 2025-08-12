package notes

import (
	"context"

	"github.com/Mafit1/notes-app/internal/models"
)

//go:generate mockgen -destination=mocks/mock_$GOFILE -package=mocks . Service
type Service interface {
	Create(ctx context.Context, note models.Note) (int64, error)
	GetAll(ctx context.Context) ([]models.Note, error)
	GetByID(ctx context.Context, id int64) (models.Note, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, note models.Note) (updatedNote models.Note, err error)
}
