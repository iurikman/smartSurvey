package server

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	//	readHeaderTimeout       = 10 * time.Second
	gracefulShutdownTimeout = 10 * time.Second
)

type metrics interface {
	TrackHttpRequest(start time.Time, req *http.Request)
}

type Server struct {
	port    string
	server  *http.Server
	key     *rsa.PublicKey
	metrics metrics
	router  *chi.Mux
	//	service service
	cfg Config
}

type Config struct {
	BindAddress string
}

func NewServer(cfg Config, key *rsa.PublicKey, metrics metrics) *Server {
	router := chi.NewRouter()

	return &Server{
		cfg: cfg,
		//		service: service,
		router:  router,
		key:     key,
		metrics: metrics,
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
