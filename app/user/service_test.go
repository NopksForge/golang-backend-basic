package user_test

import (
	"context"
	"errors"
	"testing"
	"training/app/user"
	"training/persistence"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestServiceDelete(t *testing.T) {
	tests := []struct {
		name          string
		repositoryErr error
		userId        uuid.UUID
		expectedErr   error
	}{
		{
			"1",
			nil,
			uuid.New(),
			nil,
		}, {
			"2",
			errors.New("repository error"),
			uuid.New(),
			errors.New("repository error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			repositoryMock := user.NewMockRepository(ctrl)
			repositoryMock.EXPECT().Delete(ctx, gomock.Eq(tt.userId)).Return(tt.repositoryErr).Times(1)

			s := user.NewService(repositoryMock)
			err := s.Delete(ctx, tt.userId)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestServiceGetById(t *testing.T) {
	tests := []struct {
		name          string
		user          *persistence.User
		repositoryErr error
		userId        uuid.UUID
		expectedUser  *user.GetUserDto
		expectedErr   error
	}{
		{
			"1",
			&persistence.User{},
			nil,
			uuid.New(),
			&user.GetUserDto{},
			nil,
		}, {
			"2",
			&persistence.User{},
			errors.New("repository error"),
			uuid.New(),
			nil,
			errors.New("repository error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			repositoryMock := user.NewMockRepository(ctrl)

			repositoryMock.EXPECT().SelectById(ctx, gomock.Eq(tt.userId)).Return(tt.user, tt.repositoryErr).Times(1)

			s := user.NewService(repositoryMock)
			got, err := s.GetById(ctx, tt.userId)

			assert.Equal(t, tt.expectedUser, got)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestServiceSave(t *testing.T) {
	tests := []struct {
		name          string
		user          *persistence.User
		repositoryErr error
		payload       user.SaveUserPayload
		expectedUser  *user.SaveUserDto
		expectedErr   error
	}{
		{
			"1",
			&persistence.User{},
			nil,
			user.SaveUserPayload{},
			&user.SaveUserDto{},
			nil,
		}, {
			"2",
			nil,
			errors.New("repository error"),
			user.SaveUserPayload{},
			nil,
			errors.New("repository error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			repositoryMock := user.NewMockRepository(ctrl)

			repositoryMock.EXPECT().Insert(ctx, gomock.Any()).Return(tt.user, tt.repositoryErr).Times(1)

			s := user.NewService(repositoryMock)
			got, err := s.Save(ctx, tt.payload)

			assert.Equal(t, tt.expectedUser, got)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestServiceUpdate(t *testing.T) {
	tests := []struct {
		name          string
		repositoryErr error
		payload       user.UpdateUserPayload
		expectedErr   error
	}{
		{
			"1",
			nil,
			user.UpdateUserPayload{},
			nil,
		}, {
			"2",
			errors.New("repository error"),
			user.UpdateUserPayload{},
			errors.New("repository error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			repositoryMock := user.NewMockRepository(ctrl)

			repositoryMock.EXPECT().Update(ctx, gomock.Any()).Return(tt.repositoryErr).Times(1)

			s := user.NewService(repositoryMock)
			err := s.Update(ctx, tt.payload)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
