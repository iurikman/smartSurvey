package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/iurikman/smartSurvey/internal/model"
)

type store interface {
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type Service struct {
	db store
}

func New(db store) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	rUser, err := s.db.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("s.db.CreateUser(ctx, user): %w", err)
	}

	return rUser, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.db.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("s.db.GetUserByID(ctx, id): %w", err)
	}

	return user, nil
}
