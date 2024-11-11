package user

import (
	"context"
	"time"
	"training/persistence"

	"github.com/google/uuid"
)

type Service interface {
	Save(ctx context.Context, payload SaveUserPayload) (*SaveUserDto, error)
	Update(ctx context.Context, payload UpdateUserPayload) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (*GetUserDto, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Save(ctx context.Context, payload SaveUserPayload) (*SaveUserDto, error) {
	user, err := s.repository.Insert(ctx, persistence.User{
		UserEmail: payload.UserEmail,
		UserName:  payload.UserName,
		CreatedBy: "ADMIN",
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return &SaveUserDto{
		UserId:    user.UserId,
		UserEmail: user.UserEmail,
		UserName:  user.UserName,
	}, nil
}

func (s *service) Update(ctx context.Context, payload UpdateUserPayload) error {
	user := "ADMIN"
	now := time.Now()

	return s.repository.Update(ctx, persistence.User{
		UserId:    payload.UserId,
		UserEmail: payload.UserEmail,
		UserName:  payload.UserName,
		UpdatedBy: &user,
		UpdatedAt: &now,
	})
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) GetById(ctx context.Context, id uuid.UUID) (*GetUserDto, error) {
	user, err := s.repository.SelectById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &GetUserDto{
		UserId:    user.UserId,
		UserEmail: user.UserEmail,
		UserName:  user.UserName,
	}, nil
}
