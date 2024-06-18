package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/iurikman/smartSurvey/internal/models"
)

type store interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, user models.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	CreateCompany(ctx context.Context, company models.Company) (*models.Company, error)
	UpdateCompany(ctx context.Context, company models.Company) (*models.Company, error)
}

type Service struct {
	db store
}

func New(db store) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateCompany(ctx context.Context, company models.Company) (*models.Company, error) {
	if err := company.Validate(); err != nil {
		return nil, fmt.Errorf("company.Validate(): %w", models.ErrCompanyNameIsEmpty)
	}

	rCompany, err := s.db.CreateCompany(ctx, company)
	if err != nil {
		return nil, fmt.Errorf("s.db.CreateCompany(ctx, company) err: %w", err)
	}

	return rCompany, nil
}

func (s *Service) UpdateCompany(ctx context.Context, company models.Company) (*models.Company, error) {
	if err := company.Validate(); err != nil {
		return nil, fmt.Errorf("company.Validate(): %w", err)
	}

	rCompany, err := s.db.UpdateCompany(ctx, company)
	if err != nil {
		return nil, fmt.Errorf("s.db.UpdateCompany(ctx, id) err: %w", err)
	}

	return rCompany, nil
}

func (s *Service) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("user.Validate(): %w", err)
	}

	rUser, err := s.db.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("s.db.CreateUser(ctx, user): %w", err)
	}

	return rUser, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.db.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("s.db.GetUserByID(ctx, id): %w", err)
	}

	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, user models.UpdateUserRequest) (*models.User, error) {
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("user.Validate(): %w", err)
	}

	newUser, err := s.db.UpdateUser(ctx, id, user)
	if err != nil {
		return nil, fmt.Errorf("s.db.PatchUser(ctx, user): %w", err)
	}

	return newUser, nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.db.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("s.db.DeleteUser(ctx, id): %w", err)
	}

	return nil
}
