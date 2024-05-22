package service

import (
	"context"
	"fmt"

	"github.com/iurikman/smartSurvey/internal/model"
)

type store interface {
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
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

	fmt.Println("creating user")

	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return rUser, nil
}
