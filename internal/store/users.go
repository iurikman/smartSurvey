package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/iurikman/smartSurvey/internal/model"
)

func (p *Postgres) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	query := `
			INSERT INTO users (id, name, email)
			VALUES ($1, $2, $3)
			RETURNING id, name, email
			`

	err := p.db.QueryRow(
		ctx,
		query,
		uuid.New(),
		user.Name,
		user.Email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &user, nil
}

func (p *Postgres) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := new(model.User)

	query := `
		SELECT id, name, email
		FROM users
		WHERE id = $1
		`

	err := p.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}
