package notes

import (
	"context"
	"errors"

	"github.com/Mafit1/notes-app/internal/models"
	"github.com/Mafit1/notes-app/internal/repository"
)

type Service struct {
	noteRepository NoteRepository
}

func New(repo NoteRepository) *Service {
	return &Service{repo}
}

func (s *Service) Create(ctx context.Context, note models.Note) (id int64, err error) {
	id, err = s.noteRepository.Create(ctx, note)
	if err != nil {
		return 0, ErrCannotCreateNote
	}
	return id, nil
}

func (s *Service) GetAll(ctx context.Context) (notes []models.Note, err error) {
	notes, err = s.noteRepository.GetAll(ctx)
	if err != nil {
		return nil, ErrCannotGetAllNotes
	}
	return notes, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (note *models.Note, err error) {
	note, err = s.noteRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoteNotFound) {
			return nil, ErrNoteNotFound
		}
		return nil, ErrCannotGetNote
	}
	return note, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	err := s.noteRepository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoteNotFound) {
			return ErrNoteNotFound
		}
		return ErrCannotGetNote
	}
	return nil
}

func (s *Service) Update(ctx context.Context, note models.Note) error {
	err := s.noteRepository.Update(ctx, note)
	if err != nil {
		if errors.Is(err, repository.ErrNoteNotFound) {
			return ErrNoteNotFound
		}
		return ErrCannotUpdateNote
	}
	return nil
}
