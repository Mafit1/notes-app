package notes

import (
	"context"
	"errors"

	"github.com/Mafit1/notes-app/internal/models"
	notes_repo "github.com/Mafit1/notes-app/internal/repository/notes"
)

type service struct {
	notesRepository notes_repo.Repository
}

func New(repo notes_repo.Repository) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, note models.Note) (id int64, err error) {
	id, err = s.notesRepository.Create(ctx, note)
	if err != nil {
		return 0, ErrCannotCreateNote
	}
	return id, nil
}

func (s *service) GetAll(ctx context.Context) (notes []models.Note, err error) {
	notes, err = s.notesRepository.GetAll(ctx)
	if err != nil {
		return nil, ErrCannotGetAllNotes
	}
	return notes, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (note *models.Note, err error) {
	note, err = s.notesRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, notes_repo.ErrNoteNotFound) {
			return nil, ErrNoteNotFound
		}
		return nil, ErrCannotGetNote
	}
	return note, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	err := s.notesRepository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, notes_repo.ErrNoteNotFound) {
			return ErrNoteNotFound
		}
		return ErrCannotGetNote
	}
	return nil
}

func (s *service) Update(ctx context.Context, note models.Note) error {
	err := s.notesRepository.Update(ctx, note)
	if err != nil {
		if errors.Is(err, notes_repo.ErrNoteNotFound) {
			return ErrNoteNotFound
		}
		return ErrCannotUpdateNote
	}
	return nil
}
