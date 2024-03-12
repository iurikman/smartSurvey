package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	go func() {
		http.ListenAndServe(":8000", &handler1{})
	}()
	http.ListenAndServe(":8080", &handler2{})
}

type handler1 struct {
}

func (h1 *handler1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("8000"))
	w.Write([]byte(time.Now().String()))
	s := ipStat{}
	s.printIpStat(w, r)
}

type handler2 struct {
}

func (h2 *handler2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("8080"))
	var ipAddreses []byte
	fmt.Append(ipAddreses, []byte(r.RemoteAddr))
	w.Write([]byte(r.RemoteAddr))
	s := ipStat{}
	s.printIpStat(w, r)
}

type ipStat struct {
}

var stat []string

func (i *ipStat) printIpStat(w http.ResponseWriter, r *http.Request) {

	stat = append(stat, r.RemoteAddr)
	fmt.Println(stat)
	for _, value := range stat {
		w.Write([]byte(value))
	}

}
