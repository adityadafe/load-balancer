package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type Server struct {
	Addr string
}

func NewServer(addr string) *Server {
	return &Server{
		Addr: addr,
	}
}

func (s *Server) Run() {
	port := s.Addr
	sm := http.NewServeMux()

	sm.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is healthy"))
	})

	sm.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Print("Server hosted at ", port, " served request to ", r.UserAgent())
		// converting port to int
		pti, _ := strconv.Atoi(port[1:])
		w.Write([]byte("Hello from server " + strconv.Itoa(pti-2999)))
	})

	fmt.Println("Starting server ... ")

	err := http.ListenAndServe(s.Addr, sm)

	for err != nil {
		port = fmt.Sprintf(":%d", HandlePortError(port[1:]))
		err = http.ListenAndServe(port, sm)
	}
}

func HandlePortError(pp string) int {
	// pp is previous port :)
	p, _ := strconv.Atoi(pp)
	return p + 1
}
