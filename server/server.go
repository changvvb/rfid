package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	port   int
	server string
	mux    *http.ServeMux
}

func New(port int) *Server {
	server := Server{port: port}
	server.mux = http.NewServeMux()
	server.mux.Handle("/login", http.HandlerFunc(login))
	return &server
}

func (s *Server) Run() {
	http.ListenAndServe(":"+fmt.Sprint(s.port), s.mux)
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		fmt.Fprintf(w, "hahaha")
	}
}
