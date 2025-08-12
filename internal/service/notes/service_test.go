package notes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mafit1/notes-app/internal/models"
	notes_repo "github.com/Mafit1/notes-app/internal/repository/notes"
	notes_repo_mocks "github.com/Mafit1/notes-app/internal/repository/notes/mocks"
	notes_service "github.com/Mafit1/notes-app/internal/service/notes"
	"github.com/golang/mock/gomock"
)

func TestCreate(t *testing.T) {
	var (
		ctx    = context.Background()
		noteID = int64(1)
	)

	type MockBehaivor func(r *notes_repo_mocks.MockRepository)

	note := models.Note{
		Title:   "title",
		Content: "content",
	}

	tests := []struct {
		name         string
		mockBehaivor MockBehaivor
		want         int64
		wantErr      error
	}{
		{
			name: "success",
			mockBehaivor: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().Create(ctx, note).Return(noteID, nil)
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "cannot create note",
			mockBehaivor: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().Create(ctx, note).Return(int64(0), assert.AnError)
			},
			want:    0,
			wantErr: notes_service.ErrCannotCreateNote,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockRepo := notes_repo_mocks.NewMockRepository(ctrl)
			tc.mockBehaivor(mockRepo)

			s := notes_service.New(mockRepo)

			got, err := s.Create(ctx, note)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetAll(t *testing.T) {
	var (
		ctx = context.Background()
	)

	type MockBehaivor func(r *notes_repo_mocks.MockRepository)

	notes := []models.Note{
		{
			Title:   "title 1",
			Content: "content 1",
		},
		{
			Title:   "title 2",
			Content: "content 2",
		},
		{
			Title:   "title 3",
			Content: "content 3",
		},
	}

	tests := []struct {
		name         string
		mockBehaivor MockBehaivor
		want         []models.Note
		wantErr      error
	}{
		{
			name: "success",
			mockBehaivor: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetAll(ctx).Return(notes, nil)
			},
			want:    notes,
			wantErr: nil,
		},
		{
			name: "cannot get all notes",
			mockBehaivor: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetAll(ctx).Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: notes_service.ErrCannotGetAllNotes,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockRepo := notes_repo_mocks.NewMockRepository(ctrl)
			tc.mockBehaivor(mockRepo)

			s := notes_service.New(mockRepo)

			got, err := s.GetAll(ctx)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetByID(t *testing.T) {
	var (
		ctx       = context.Background()
		id        = int64(1)
		emptyNote = models.Note{}
	)

	type MockBehaivor func(r *notes_repo_mocks.MockRepository)

	note := models.Note{
		Title:   "title",
		Content: "content",
	}

	tests := []struct {
		name         string
		mockBehaivor MockBehaivor
		want         models.Note
		wantErr      error
	}{
		{
			name: "success",
			mockBehaivor: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, id).Return(note, nil)
			},
			want:    note,
			wantErr: nil,
		},
		{
			name: "note not found",
			mockBehaivor: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, id).Return(emptyNote, notes_repo.ErrNoteNotFound)
			},
			want:    emptyNote,
			wantErr: notes_service.ErrNoteNotFound,
		},
		{
			name: "database error",
			mockBehaivor: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, id).Return(emptyNote, notes_repo.ErrDatabase)
			},
			want:    emptyNote,
			wantErr: notes_service.ErrCannotGetNote,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockRepo := notes_repo_mocks.NewMockRepository(ctrl)
			tc.mockBehaivor(mockRepo)

			s := notes_service.New(mockRepo)

			got, err := s.GetByID(ctx, id)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDelete(t *testing.T) {
	var (
		ctx = context.Background()
		id  = int64(1)
	)

	type MockBehaivor func(r *notes_repo_mocks.MockRepository)

	tests := []struct {
		name         string
		mockBehaivor MockBehaivor
		wantErr      error
	}{
		{
			name: "success",
			mockBehaivor: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().Delete(ctx, id).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "note not found",
			mockBehaivor: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().Delete(ctx, id).Return(notes_repo.ErrNoteNotFound)
			},
			wantErr: notes_service.ErrNoteNotFound,
		},
		{
			name: "database error",
			mockBehaivor: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().Delete(ctx, id).Return(notes_repo.ErrDatabase)
			},
			wantErr: notes_service.ErrCannotDeleteNote,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockRepo := notes_repo_mocks.NewMockRepository(ctrl)
			tc.mockBehaivor(mockRepo)

			s := notes_service.New(mockRepo)

			err := s.Delete(ctx, id)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestUpdate(t *testing.T) {
	var (
		ctx       = context.Background()
		noteID    = int64(1)
		emptyNote = models.Note{}
	)

	type MockBehaivor func(r *notes_repo_mocks.MockRepository)

	note := models.Note{
		ID:      noteID,
		Title:   "title",
		Content: "content",
	}

	updatedNote := models.Note{
		ID:      noteID,
		Title:   "updated title",
		Content: "updated content",
	}

	tests := []struct {
		name         string
		mockBehaivor MockBehaivor
		want         models.Note
		wantErr      error
	}{
		{
			name: "success",
			mockBehaivor: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().Update(ctx, note).Return(updatedNote, nil)
			},
			want:    updatedNote,
			wantErr: nil,
		},
		{
			name: "note not found",
			mockBehaivor: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().Update(ctx, note).Return(emptyNote, notes_repo.ErrNoteNotFound)
			},
			want:    emptyNote,
			wantErr: notes_service.ErrNoteNotFound,
		},
		{
			name: "database error",
			mockBehaivor: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().Update(ctx, note).Return(emptyNote, notes_repo.ErrDatabase)
			},
			want:    emptyNote,
			wantErr: notes_service.ErrCannotUpdateNote,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockRepo := notes_repo_mocks.NewMockRepository(ctrl)
			tc.mockBehaivor(mockRepo)

			s := notes_service.New(mockRepo)

			got, err := s.Update(ctx, note)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}
