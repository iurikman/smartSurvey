package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iurikman/smartSurvey/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *Postgres) CreateCompany(ctx context.Context, company models.Company) (*models.Company, error) {
	query := `
			INSERT INTO companies (id, name)
			VALUES ($1, $2)
			RETURNING id, name 
	`

	err := p.db.QueryRow(
		ctx,
		query,
		uuid.New(),
		company.Name,
	).Scan(
		&company.ID,
		&company.Name,
	)

	var pgErr *pgconn.PgError

	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, models.ErrDuplicateCompany
		}

		return nil, fmt.Errorf("error creating company: %w", err)
	}

	return &company, nil
}

func (p *Postgres) UpdateCompany(ctx context.Context, company models.Company) (*models.Company, error) {
	var changedCompany models.Company

	query := `
		UPDATE companies
		SET name = $2
		WHERE id = $1
		RETURNING id, name
	`

	err := p.db.QueryRow(
		ctx,
		query,
		company.ID,
		company.Name,
	).Scan(
		&changedCompany.ID,
		&changedCompany.Name,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, models.ErrCompanyNotFound
	case err != nil:
		return nil, fmt.Errorf("error updating company: %w", err)
	}

	return &company, nil
}
