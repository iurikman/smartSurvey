package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func main() {
	r := mux.NewRouter()
	var h handler
	r.HandleFunc("/time", h.handleTime)
	r.HandleFunc("/stats", h.handleStats)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	err := srv.ListenAndServe()
	if err != nil {
		return
	}
}

type handler struct {
	ipStats ipStats
}

func (h *handler) handleStats(w http.ResponseWriter, r *http.Request) {
	h.ipStats.ipInfo = make(map[string]int)
	h.ipStats.printIpStats(w, r)
}

func (h *handler) handleTime(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(time.Now().String()))
	if err != nil {
		return
	}
}

type ipStats struct {
	ipInfo map[string]int
}

func (i *ipStats) printIpStats(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	var ipStatInString string
	_, ok := i.ipInfo[ip]
	if ok {
		i.ipInfo[ip]++
	} else {
		i.ipInfo[ip] = 1
	}

	for key, val := range i.ipInfo {
		ipStatInString += fmt.Sprint(key + " :  " + strconv.Itoa(val) + "  ||||  ")
	}
	_, err := w.Write([]byte(ipStatInString))
	if err != nil {
		return
	}
}
