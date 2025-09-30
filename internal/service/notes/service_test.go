package notes_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Mafit1/notes-app/internal/metrics/mocks"
	notes_repo "github.com/Mafit1/notes-app/internal/repository/notes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Mafit1/notes-app/internal/models"
	notes_repo_mocks "github.com/Mafit1/notes-app/internal/repository/notes/mocks"
	notes_service "github.com/Mafit1/notes-app/internal/service/notes"
	"github.com/golang/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	note := notes_service.CreateNote{Title: "title", Content: "content"}

	type MockBehavior func(r *notes_repo_mocks.MockRepository, m *mocks.MockNotesMetrics)

	tests := []struct {
		name    string
		mock    MockBehavior
		wantID  int64
		wantErr error
	}{
		{
			name: "success",
			mock: func(r *notes_repo_mocks.MockRepository, m *mocks.MockNotesMetrics) {
				r.EXPECT().
					CreateByUserID(ctx, userID, notes_repo.CreateNote{Title: note.Title, Content: note.Content}).
					Return(int64(1), nil)
				m.EXPECT().
					IncCreated().
					Times(1)
			},
			wantID:  1,
			wantErr: nil,
		},
		{
			name: "cannot create note",
			mock: func(r *notes_repo_mocks.MockRepository, m *mocks.MockNotesMetrics) {
				r.EXPECT().
					CreateByUserID(ctx, userID, notes_repo.CreateNote{Title: note.Title, Content: note.Content}).
					Return(int64(0), errors.New("db error"))
				m.EXPECT().
					IncCreatedError().
					Times(1)
			},
			wantID:  0,
			wantErr: notes_service.ErrCannotCreateNote,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := notes_repo_mocks.NewMockRepository(ctrl)
			mockMetrics := mocks.NewMockNotesMetrics(ctrl)
			tc.mock(mockRepo, mockMetrics)

			s := notes_service.New(mockRepo, mockMetrics)
			id, err := s.Create(ctx, userID, note)

			assert.Equal(t, tc.wantID, id)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestGetAll(t *testing.T) {
	var (
		ctx = context.Background()
	)

	type MockBehavior func(r *notes_repo_mocks.MockRepository)

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
		mockBehavior MockBehavior
		want         []models.Note
		wantErr      error
	}{
		{
			name: "success",
			mockBehavior: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetAll(ctx).Return(notes, nil)
			},
			want:    notes,
			wantErr: nil,
		},
		{
			name: "cannot get all notes",
			mockBehavior: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetAll(ctx).Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: notes_service.ErrCannotGetNotes,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			mockRepo := notes_repo_mocks.NewMockRepository(ctrl)
			mockMetrics := mocks.NewMockNotesMetrics(ctrl)
			tc.mockBehavior(mockRepo)

			s := notes_service.New(mockRepo, mockMetrics)

			got, err := s.GetAll(ctx)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	note := models.Note{ID: 1, Title: "title", Content: "content", UserID: userID}

	type MockBehavior func(r *notes_repo_mocks.MockRepository)

	tests := []struct {
		name    string
		mock    MockBehavior
		want    *models.Note
		wantErr error
	}{
		{
			name: "success",
			mock: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, int64(1)).Return(note, nil)
			},
			want:    &note,
			wantErr: nil,
		},
		{
			name: "note not found",
			mock: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, int64(1)).Return(models.Note{}, notes_repo.ErrNoteNotFound)
			},
			want:    nil,
			wantErr: notes_service.ErrNoteNotFound,
		},
		{
			name: "cannot get note (db error)",
			mock: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, int64(1)).Return(models.Note{}, errors.New("db error"))
			},
			want:    nil,
			wantErr: notes_service.ErrCannotGetNote,
		},
		{
			name: "forbidden",
			mock: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, int64(1)).Return(models.Note{ID: 1, UserID: otherUserID}, nil)
			},
			want:    nil,
			wantErr: notes_service.ErrForbidden,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := notes_repo_mocks.NewMockRepository(ctrl)
			mockMetrics := mocks.NewMockNotesMetrics(ctrl)
			tc.mock(mockRepo)

			s := notes_service.New(mockRepo, mockMetrics)
			got, err := s.GetByID(ctx, userID, 1)

			assert.Equal(t, tc.want, got)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestGetAllByUserID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notes := []*models.Note{
		{ID: 1, Title: "title1", Content: "content1", UserID: userID},
		{ID: 2, Title: "title2", Content: "content2", UserID: userID},
	}

	type MockBehavior func(r *notes_repo_mocks.MockRepository)

	tests := []struct {
		name    string
		mock    MockBehavior
		want    []*models.Note
		wantErr error
	}{
		{
			name: "success",
			mock: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetAllFromUserByID(ctx, userID).Return(notes, nil)
			},
			want:    notes,
			wantErr: nil,
		},
		{
			name: "cannot get notes",
			mock: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetAllFromUserByID(ctx, userID).Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: notes_service.ErrCannotGetNotes,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := notes_repo_mocks.NewMockRepository(ctrl)
			mockMetrics := mocks.NewMockNotesMetrics(ctrl)
			tc.mock(mockRepo)

			s := notes_service.New(mockRepo, mockMetrics)
			got, err := s.GetAllByUserID(ctx, userID)

			assert.Equal(t, tc.want, got)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	note := models.Note{ID: 1, UserID: userID}

	type MockBehavior func(r *notes_repo_mocks.MockRepository)

	tests := []struct {
		name    string
		mock    MockBehavior
		wantErr error
	}{
		{
			name: "success",
			mock: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, int64(1)).Return(note, nil)
				r.EXPECT().Delete(ctx, int64(1)).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "note not found",
			mock: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, int64(1)).Return(models.Note{}, notes_repo.ErrNoteNotFound)
			},
			wantErr: notes_service.ErrNoteNotFound,
		},
		{
			name: "forbidden",
			mock: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, int64(1)).Return(models.Note{ID: 1, UserID: otherUserID}, nil)
			},
			wantErr: notes_service.ErrForbidden,
		},
		{
			name: "cannot delete note",
			mock: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, int64(1)).Return(note, nil)
				r.EXPECT().Delete(ctx, int64(1)).Return(errors.New("db error"))
			},
			wantErr: notes_service.ErrCannotDeleteNote,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := notes_repo_mocks.NewMockRepository(ctrl)
			mockMetrics := mocks.NewMockNotesMetrics(ctrl)
			tc.mock(mockRepo)

			s := notes_service.New(mockRepo, mockMetrics)
			err := s.Delete(ctx, userID, 1)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	note := models.Note{ID: 1, Title: "title", Content: "content", UserID: userID}

	type MockBehavior func(r *notes_repo_mocks.MockRepository)

	tests := []struct {
		name    string
		mock    MockBehavior
		want    *models.Note
		wantErr error
	}{
		{
			name: "success",
			mock: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, int64(1)).Return(note, nil)
				r.EXPECT().Update(ctx, note).Return(note, nil)
			},
			want:    &note,
			wantErr: nil,
		},
		{
			name: "note not found",
			mock: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, int64(1)).Return(models.Note{}, notes_repo.ErrNoteNotFound)
			},
			want:    nil,
			wantErr: notes_service.ErrNoteNotFound,
		},
		{
			name: "forbidden",
			mock: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, int64(1)).Return(models.Note{ID: 1, UserID: otherUserID}, nil)
			},
			want:    nil,
			wantErr: notes_service.ErrForbidden,
		},
		{
			name: "cannot update note",
			mock: func(r *notes_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, int64(1)).Return(note, nil)
				r.EXPECT().Update(ctx, note).Return(models.Note{}, errors.New("db error"))
			},
			want:    nil,
			wantErr: notes_service.ErrCannotUpdateNote,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := notes_repo_mocks.NewMockRepository(ctrl)
			mockMetrics := mocks.NewMockNotesMetrics(ctrl)
			tc.mock(mockRepo)

			s := notes_service.New(mockRepo, mockMetrics)
			got, err := s.Update(ctx, userID, note)

			assert.Equal(t, tc.want, got)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
