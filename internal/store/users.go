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

func (p *Postgres) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	query := `
			INSERT INTO users (id, company, role, name, surname, phone, email, user_type)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, company, role, name, surname, phone, email, user_type
			`

	err := p.db.QueryRow(
		ctx,
		query,
		uuid.New(),
		user.Company,
		user.Role,
		user.Name,
		user.Surname,
		user.Phone,
		user.Email,
		user.UserType,
	).Scan(
		&user.ID,
		&user.Company,
		&user.Role,
		&user.Name,
		&user.Surname,
		&user.Phone,
		&user.Email,
		&user.UserType,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, models.ErrDuplicateUser
		}

		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &user, nil
}

func (p *Postgres) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := new(models.User)

	query := `
		SELECT id, company, role, name, surname, phone, email, user_type
		FROM users
		WHERE id = $1
		`

	err := p.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Company,
		&user.Role,
		&user.Name,
		&user.Surname,
		&user.Phone,
		&user.Email,
		&user.UserType,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, models.ErrUserNotFound
	case err != nil:
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

func (p *Postgres) UpdateUser(ctx context.Context, id uuid.UUID, user models.UpdateUserRequest) (*models.User, error) {
	var changedUser models.User

	query := `
		UPDATE users
		SET company	= $2, role = $3, name = $4, surname = $5, phone = $6, email = $7, user_type = $8
		WHERE id = $1
		RETURNING company, role, name, surname, phone, email, user_type
		`

	err := p.db.QueryRow(
		ctx,
		query,
		id,
		user.Company,
		user.Role,
		user.Name,
		user.Surname,
		user.Phone,
		user.Email,
		user.UserType,
	).Scan(
		&changedUser.Company,
		&changedUser.Role,
		&changedUser.Name,
		&changedUser.Surname,
		&changedUser.Phone,
		&changedUser.Email,
		&changedUser.UserType,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, models.ErrUserNotFound
	case err != nil:
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return &changedUser, nil
}

func (p *Postgres) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
		`

	result, err := p.db.Exec(
		ctx,
		query,
		id,
	)

	if result.RowsAffected() == 0 {
		return models.ErrUserNotFound
	}

	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}
