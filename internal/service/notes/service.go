package notes

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mafit1/notes-app/internal/metrics"
	"github.com/Mafit1/notes-app/internal/models"
	notes_repo "github.com/Mafit1/notes-app/internal/repository/notes"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type service struct {
	notesRepository notes_repo.Repository
	metrics         metrics.NotesMetrics
}

func New(repo notes_repo.Repository, metrics metrics.NotesMetrics) Service {
	return &service{repo, metrics}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, note CreateNote) (id int64, err error) {
	logrus.Infof("Service: Creating note for user with ID: %s", userID)
	id, err = s.notesRepository.CreateByUserID(
		ctx,
		userID,
		notes_repo.CreateNote{
			Title:   note.Title,
			Content: note.Content,
		},
	)
	if err != nil {
		s.metrics.IncCreatedError()
		return 0, ErrCannotCreateNote
	}
	logrus.Infof("Service: Note created, ID: %d", id)
	s.metrics.IncCreated()
	return id, nil
}

func (s *service) GetAll(ctx context.Context) (notes []models.Note, err error) {
	notes, err = s.notesRepository.GetAll(ctx)
	if err != nil {
		return nil, ErrCannotGetNotes
	}
	return notes, nil
}

func (s *service) GetByID(ctx context.Context, userID uuid.UUID, noteID int64) (*models.Note, error) {
	note, err := s.notesRepository.GetByID(ctx, noteID)
	if err != nil {
		if errors.Is(err, notes_repo.ErrNoteNotFound) {
			return nil, ErrNoteNotFound
		}
		return nil, ErrCannotGetNote
	}

	if note.UserID != userID {
		return nil, ErrForbidden
	}

	return &note, nil
}

func (s *service) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Note, error) {
	notes, err := s.notesRepository.GetAllFromUserByID(ctx, userID)
	if err != nil {
		return nil, ErrCannotGetNotes
	}
	return notes, nil
}

func (s *service) Delete(ctx context.Context, userID uuid.UUID, noteID int64) error {
	note, err := s.notesRepository.GetByID(ctx, noteID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNoteNotFound, err)
	}

	if note.UserID != userID {
		return ErrForbidden
	}

	err = s.notesRepository.Delete(ctx, noteID)
	if err != nil {
		if errors.Is(err, notes_repo.ErrNoteNotFound) {
			return ErrNoteNotFound
		}
		return ErrCannotDeleteNote
	}
	return nil
}

func (s *service) Update(ctx context.Context, userID uuid.UUID, note models.Note) (*models.Note, error) {
	n, err := s.notesRepository.GetByID(ctx, note.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoteNotFound, err)
	}

	if n.UserID != userID {
		return nil, ErrForbidden
	}

	updatedNote, err := s.notesRepository.Update(ctx, note)
	if err != nil {
		if errors.Is(err, notes_repo.ErrNoteNotFound) {
			return nil, ErrNoteNotFound
		}
		return nil, ErrCannotUpdateNote
	}
	return &updatedNote, nil
}
