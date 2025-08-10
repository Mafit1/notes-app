package notes_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mafit1/notes-app/internal/models"
	notes_repo_mocks "github.com/Mafit1/notes-app/internal/repository/notes/mocks"
	notes_service "github.com/Mafit1/notes-app/internal/service/notes"
	"github.com/golang/mock/gomock"
)

func TestCreate(t *testing.T) {

}

func TestGetAll(t *testing.T) {
	var (
		arbitraryErr = errors.New("arbitrary error")
		ctx          = context.Background()
	)

	type MockBehaivor func(s *notes_repo_mocks.MockRepository)

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
				r.EXPECT().GetAll(ctx).Return(nil, arbitraryErr)
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

}

func TestDelete(t *testing.T) {

}

func TestUpdate(t *testing.T) {

}
