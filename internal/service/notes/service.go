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
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"component": "note_service",
		"user_id":   userID,
		"op":        "Create",
	})

	logger.Debug("start creating note")

	id, err = s.notesRepository.CreateByUserID(ctx, userID, notes_repo.CreateNote{
		Title:   note.Title,
		Content: note.Content,
	})
	if err != nil {
		logger.WithError(err).Error("failed to create note")
		s.metrics.IncCreatedError()
		return 0, ErrCannotCreateNote
	}

	logger.WithField("note_id", id).Info("note created")
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
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"component": "note_service",
		"user_id":   userID,
		"op":        "GetByID",
	})

	logger.Debug("start getting note by id")

	note, err := s.notesRepository.GetByID(ctx, noteID)
	if err != nil {
		if errors.Is(err, notes_repo.ErrNoteNotFound) {
			logger.WithError(err).Warn("note not found")
			return nil, ErrNoteNotFound
		}
		logger.WithError(err).Error("failed to get note")
		return nil, ErrCannotGetNote
	}

	if note.UserID != userID {
		logger.Warn("access denied: note doesn't belong to user")
		return nil, ErrForbidden
	}

	logger.WithField("note_id", note.ID).Info("note received successfully")
	return &note, nil
}

func (s *service) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Note, error) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"component": "note_service",
		"user_id":   userID,
		"op":        "GetAllByUserID",
	})

	logger.Debug("start getting all notes by user id")

	notes, err := s.notesRepository.GetAllFromUserByID(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get all notes from user")
		return nil, ErrCannotGetNotes
	}

	logger.Info("notes received successfully")
	return notes, nil
}

func (s *service) Delete(ctx context.Context, userID uuid.UUID, noteID int64) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"component": "note_service",
		"user_id":   userID,
		"note_id":   noteID,
		"op":        "Delete",
	})

	logger.Debug("start deleting note")

	note, err := s.notesRepository.GetByID(ctx, noteID)
	if err != nil {
		if errors.Is(err, notes_repo.ErrNoteNotFound) {
			logger.WithError(err).Warn("note not found during deletion")
			return ErrNoteNotFound
		}
		logger.WithError(err).Error("failed to get note for deletion")
		return fmt.Errorf("get note: %w", err)
	}

	if note.UserID != userID {
		logger.Warn("access denied: note doesn't belong to user")
		return ErrForbidden
	}

	err = s.notesRepository.Delete(ctx, noteID)
	if err != nil {
		logger.WithError(err).Error("failed to delete note")
		return ErrCannotDeleteNote
	}

	logger.Info("note deleted successfully")
	return nil
}

func (s *service) Update(ctx context.Context, userID uuid.UUID, note models.Note) (*models.Note, error) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"component": "note_service",
		"user_id":   userID,
		"note_id":   note.ID,
		"op":        "Update",
	})

	logger.Debug("start updating note")

	n, err := s.notesRepository.GetByID(ctx, note.ID)
	if err != nil {
		if errors.Is(err, notes_repo.ErrNoteNotFound) {
			logger.WithError(err).Warn("note not found during update")
			return nil, ErrNoteNotFound
		}
		logger.WithError(err).Error("failed to get note for update")
		return nil, fmt.Errorf("get note: %w", err)
	}

	if n.UserID != userID {
		return nil, ErrForbidden
	}

	updatedNote, err := s.notesRepository.Update(ctx, note)
	if err != nil {
		logger.WithError(err).Error("failed to update note")
		return nil, ErrCannotUpdateNote
	}

	logger.Info("note updated")
	return &updatedNote, nil
}
