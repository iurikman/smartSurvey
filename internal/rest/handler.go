package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/iurikman/smartSurvey/internal/model"
)

type HTTPResponse struct {
	Data any   `json:"data"`
	Err  error `json:"error"`
}

type TransferResponse struct {
	TransactionID uuid.UUID `json:"transactionId"`
}

type service interface {
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
}

func (s *Server) createUser(name string) {
	var rUser model.User
	rUser.Name = name
	user, _ := s.service.CreateUser(context.Background(), rUser)
}
