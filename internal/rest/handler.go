package server

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/google/uuid"
	"github.com/iurikman/smartSurvey/internal/model"
	log "github.com/sirupsen/logrus"
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
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var rUser model.User

	user, err := s.service.CreateUser(r.Context(), rUser)
	if err != nil {
		log.Warnf("s.service.CreateUser err: %v", err)
	}

	writeOkResponse(w, http.StatusCreated, user)
}

func (s *Server) getUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		log.Warnf("uuid.Parse err: %v", err)

		return
	}

	user, err := s.service.GetUserByID(r.Context(), id)
	if err != nil {
		log.Warnf("s.service.GetUserByID err: %v", err)

		return
	}
	writeOkResponse(w, http.StatusOK, user)

}

func writeOkResponse(w http.ResponseWriter, statusCode int, user any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(HTTPResponse{Data: user})
	if err != nil {
		log.Warn("writeOkResponse/json.NewEncoder(w).Encode(HTTPResponse{Data: data})")
	}
}
