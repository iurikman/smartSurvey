package model

import "github.com/google/uuid"

type ctxKey string

const (
	UserInfoKey  ctxKey = "userInfo"
	Standartpage int    = 10
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Company  uuid.UUID `json:"company"`
	Role     uuid.UUID `json:"role"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Phone    int       `json:"phone"`
	Email    string    `json:"email"`
	UserType string    `json:"userType"`
}

type Company struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
