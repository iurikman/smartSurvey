package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
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
		ReadHeaderTimeout: 10 * time.Second,
		Addr:              port,
		Handler:           r,
	}

	return &Server{port: srv.Addr, server: srv}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		ctxWithTimeOut, cancel := context.WithTimeout(ctx, 10*time.Second)

		defer cancel()

		err := s.server.Shutdown(ctxWithTimeOut)
		if err != nil {
			log.Panic("server Shutdown error")
		}
	}()

	err := s.server.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) && err != nil {
		log.Panic("ListenAndServe error")

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
		log.Panic("Write error")

		return
	}
}

func (h *handler) handleTime(w http.ResponseWriter, r *http.Request) {
	h.ipStats.ipInfo[r.RemoteAddr]++

	_, err := w.Write([]byte(time.Now().String()))
	if err != nil {
		log.Panic("Write error")

		return
	}
}

type ipStats struct {
	ipInfo map[string]int
}
