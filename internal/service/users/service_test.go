package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Mafit1/notes-app/internal/models"
	users_repo "github.com/Mafit1/notes-app/internal/repository/users"
	users_repo_mocks "github.com/Mafit1/notes-app/internal/repository/users/mocks"
	users_service "github.com/Mafit1/notes-app/internal/service/users"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	createReq := users_service.CreateUser{
		Name:           "John",
		Email:          "john@example.com",
		HashedPassword: "hashed-pass",
	}

	expectedRepoReq := users_repo.CreateUser{
		Name:     "John",
		Email:    "john@example.com",
		Password: "hashed-pass",
	}

	successID := uuid.New()

	type MockBehavior func(r *users_repo_mocks.MockRepository)

	tests := []struct {
		name    string
		mock    MockBehavior
		wantID  uuid.UUID
		wantErr error
	}{
		{
			name: "success",
			mock: func(r *users_repo_mocks.MockRepository) {
				r.EXPECT().Create(ctx, expectedRepoReq).Return(successID, nil)
			},
			wantID:  successID,
			wantErr: nil,
		},
		{
			name: "cannot create user",
			mock: func(r *users_repo_mocks.MockRepository) {
				r.EXPECT().Create(ctx, expectedRepoReq).Return(uuid.Nil, errors.New("db error"))
			},
			wantID:  uuid.Nil,
			wantErr: users_service.ErrCannotCreateUser,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := users_repo_mocks.NewMockRepository(ctrl)
			tc.mock(mockRepo)

			s := users_service.New(mockRepo)
			gotID, err := s.Create(ctx, createReq)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantID, gotID)
		})
	}
}

func TestGetByEmail(t *testing.T) {
	ctx := context.Background()
	email := "john@example.com"
	user := &models.User{ID: uuid.New(), Name: "John", Email: email}

	type MockBehavior func(r *users_repo_mocks.MockRepository)

	tests := []struct {
		name    string
		mock    MockBehavior
		want    *models.User
		wantErr error
	}{
		{
			name: "success",
			mock: func(r *users_repo_mocks.MockRepository) {
				r.EXPECT().GetByEmail(ctx, email).Return(user, nil)
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "user not found",
			mock: func(r *users_repo_mocks.MockRepository) {
				r.EXPECT().GetByEmail(ctx, email).Return(nil, users_repo.ErrUserNotFound)
			},
			want:    nil,
			wantErr: users_service.ErrUserNotFound,
		},
		{
			name: "cannot get user",
			mock: func(r *users_repo_mocks.MockRepository) {
				r.EXPECT().GetByEmail(ctx, email).Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: users_service.ErrCannotGetUser,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := users_repo_mocks.NewMockRepository(ctrl)
			tc.mock(mockRepo)

			s := users_service.New(mockRepo)
			got, err := s.GetByEmail(ctx, email)

			assert.Equal(t, tc.want, got)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	user := &models.User{ID: userID, Name: "John", Email: "john@example.com"}

	type MockBehavior func(r *users_repo_mocks.MockRepository)

	tests := []struct {
		name    string
		mock    MockBehavior
		want    *models.User
		wantErr error
	}{
		{
			name: "success",
			mock: func(r *users_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, userID).Return(user, nil)
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "user not found",
			mock: func(r *users_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, userID).Return(nil, users_repo.ErrUserNotFound)
			},
			want:    nil,
			wantErr: users_service.ErrUserNotFound,
		},
		{
			name: "cannot get user",
			mock: func(r *users_repo_mocks.MockRepository) {
				r.EXPECT().GetByID(ctx, userID).Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: users_service.ErrCannotGetUser,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := users_repo_mocks.NewMockRepository(ctrl)
			tc.mock(mockRepo)

			s := users_service.New(mockRepo)
			got, err := s.GetByID(ctx, userID)

			assert.Equal(t, tc.want, got)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
