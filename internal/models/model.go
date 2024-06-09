package models

import "github.com/google/uuid"

type UpdateUserRequest struct {
	Company  *uuid.UUID `json:"company"`
	Role     string     `json:"role"`
	Name     string     `json:"name"`
	Surname  string     `json:"surname"`
	Phone    string     `json:"phone"`
	Email    string     `json:"email"`
	UserType string     `json:"userType"`
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

type Company struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
