package server

import (
	"context"
	"net/http"

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

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var rUser model.User

	_, _ = s.service.CreateUser(r.Context(), rUser)
}
