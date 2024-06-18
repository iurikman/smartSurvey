package models

import (
	"github.com/google/uuid"
)

type UpdateUserRequest struct {
	Company  *uuid.UUID `json:"company"`
	Role     string     `json:"role"`
	Name     string     `json:"name"`
	Surname  string     `json:"surname"`
	Phone    string     `json:"phone"`
	Email    string     `json:"email"`
	UserType string     `json:"userType"`
}

func (u UpdateUserRequest) Validate() error {
	if u.Name == "" {
		return ErrUserNameIsEmpty
	}

	if u.Email == "" {
		return ErrEmailIsEmpty
	}

	if u.Phone == "" {
		return ErrPhoneIsEmpty
	}

	return nil
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Company  uuid.UUID `json:"company"`
	Role     string    `json:"role"`
	Name     string    `json:"name"`
	Surname  *string   `json:"surname"`
	Phone    *string   `json:"phone"`
	Email    *string   `json:"email"`
	UserType *string   `json:"userType"`
}

func (u User) Validate() error {
	if u.Name == "" {
		return ErrUserNameIsEmpty
	}

	if *u.Email == "" {
		return ErrEmailIsEmpty
	}

	if *u.Phone == "" {
		return ErrPhoneIsEmpty
	}

	return nil
}

type Company struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (c Company) Validate() error {
	if c.Name == "" {
		return ErrCompanyNameIsEmpty
	}

	return nil
}
