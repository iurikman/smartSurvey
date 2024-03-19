package server

import (
	"fmt"
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
		Addr:    ":8080",
		Handler: r,
	}
	return &Server{port: srv.Addr, server: srv}
}

func (s *Server) Start() {

	err := s.server.ListenAndServe()
	if err != http.ErrServerClosed && err != nil {
		panic("srever error")
	}
}

type handler struct {
	ipStats ipStats
}

func (h *handler) handleStats(w http.ResponseWriter, r *http.Request) {
	var ipStatInString string

	for key, val := range h.ipStats.ipInfo {
		ipStatInString += fmt.Sprint(key + " :  " + strconv.Itoa(val) + "  ||||  ")
	}
	_, err := w.Write([]byte(ipStatInString))
	if err != nil {
		return
	}
}

func (h *handler) handleTime(w http.ResponseWriter, r *http.Request) {
	h.ipStats.ipInfo[r.RemoteAddr]++

	_, err := w.Write([]byte(time.Now().String()))
	if err != nil {
		return
	}
}

type ipStats struct {
	ipInfo map[string]int
}
