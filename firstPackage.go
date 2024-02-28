package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type TimeService struct {
}

func (s *TimeService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	fmt.Fprintf(w, "%s\n", currentTime.String())
}

func main() {
	s := &TimeService{}
	http.HandleFunc("/time", s.ServeHTTP)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
