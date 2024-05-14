package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

const (
	//	readHeaderTimeout       = 10 * time.Second
	gracefulShutdownTimeout = 10 * time.Second
)

type Server struct {
	router *chi.Mux
	cfg    Config
	server *http.Server
}

type Config struct {
	BindAddress string
}

func NewServer(cfg Config) *Server {
	router := chi.NewRouter()

	return &Server{
		cfg:    cfg,
		router: router,
		server: &http.Server{
			Addr:              cfg.BindAddress,
			ReadHeaderTimeout: 5 * time.Second,
			Handler:           router,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		ctxWithTimeOut, cancel := context.WithTimeout(ctx, gracefulShutdownTimeout)

		defer cancel()

		err := s.server.Shutdown(ctxWithTimeOut)
		if err != nil {
			logrus.Warnf("server Shutdown error: %v", err)
		}
	}()

	err := s.server.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server closed error: %w", err)
	}

	return nil
}
