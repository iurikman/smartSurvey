package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	go func() {
		err := http.ListenAndServe(":8000", &handler1{})
		if err != nil {
			return
		}
	}()
	err := http.ListenAndServe(":8080", &handler2{})
	if err != nil {
		return
	}
}

type handler1 struct{}

func (h1 *handler1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("--- 8 0 0 0 ---"))
	if err != nil {
		return
	}
	_, err = w.Write([]byte(time.Now().String()))
	if err != nil {
		return
	}
	ipStat := ipStat{}
	ipStat.printIpStat(w, r)
}

type handler2 struct{}

func (h2 *handler2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("--- 8 0 8 0 ---"))
	if err != nil {
		return
	}
}

type ipStat struct {
	ipInfo map[string]int
}

func (i *ipStat) printIpStat(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	val, ok := i.ipInfo[ip]
	if ok {
		i.ipInfo[ip] = val + 1
	} else {
		i.ipInfo[ip] = 1
	}
	ipStatInString := ""
	for key, val := range i.ipInfo {
		ipStatInString += fmt.Sprint(key + ": " + string(rune(val)) + "\n")
	}
	_, err := w.Write([]byte(ipStatInString))
	if err != nil {
		return
	}
}
