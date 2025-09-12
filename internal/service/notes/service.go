package notes

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mafit1/notes-app/internal/models"
	notes_repo "github.com/Mafit1/notes-app/internal/repository/notes"
	"github.com/google/uuid"
)

type service struct {
	notesRepository notes_repo.Repository
}

func New(repo notes_repo.Repository) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, note CreateNote) (id int64, err error) {
	id, err = s.notesRepository.Create(
		ctx,
		notes_repo.CreateNote{
			Title:   note.Title,
			Content: note.Content,
		},
	)
	if err != nil {
		return 0, ErrCannotCreateNote
	}
	return id, nil
}

func (s *service) CreateByUserID(ctx context.Context, userID uuid.UUID, note CreateNote) (id int64, err error) {
	id, err = s.notesRepository.CreateByUserID(
		ctx,
		userID,
		notes_repo.CreateNote{
			Title:   note.Title,
			Content: note.Content,
		},
	)
	if err != nil {
		return 0, ErrCannotCreateNote
	}
	return id, nil
}

func (s *service) GetAll(ctx context.Context) (notes []models.Note, err error) {
	notes, err = s.notesRepository.GetAll(ctx)
	if err != nil {
		return nil, ErrCannotGetNotes
	}
	return notes, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (note models.Note, err error) {
	note, err = s.notesRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, notes_repo.ErrNoteNotFound) {
			return models.Note{}, ErrNoteNotFound
		}
		return models.Note{}, ErrCannotGetNote
	}
	return note, nil
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
