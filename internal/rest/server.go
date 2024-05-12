package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	readHeaderTimeout       = 10 * time.Second
	gracefulShutdownTimeout = 10 * time.Second
)

type Server struct {
	port   string
	server *http.Server
}

func NewServer(port string) *Server {
	r := http.NewServeMux()
	h := handler{
		ipStats: ipStats{
			ipInfo: make(map[string]int),
		},
	}
	r.HandleFunc("GET /time", h.handleTime)
	r.HandleFunc("GET /stats", h.handleStats)
	srv := &http.Server{
		ReadHeaderTimeout: readHeaderTimeout,
		Addr:              port,
		Handler:           r,
	}

	return &Server{port: srv.Addr, server: srv}
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

type handler struct {
	ipStats ipStats
}

func (h *handler) handleStats(w http.ResponseWriter, _ *http.Request) {
	var ipStatInString string

	for key, val := range h.ipStats.ipInfo {
		ipStatInString += key + " :  " + strconv.Itoa(val) + "  ||||  "
	}

	_, err := w.Write([]byte(ipStatInString))
	if err != nil {
		logrus.Warnf("Write error: %v", err)

		return
	}
}

func (h *handler) handleTime(w http.ResponseWriter, r *http.Request) {
	h.ipStats.mx.Lock()
	defer h.ipStats.mx.Unlock()

	h.ipStats.ipInfo[r.RemoteAddr]++

	_, err := w.Write([]byte(time.Now().String()))
	if err != nil {
		logrus.Warnf("Write error: %v", err)

		return
	}
}

type ipStats struct {
	ipInfo map[string]int
	mx     sync.Mutex
}
