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
