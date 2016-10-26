package rfid

import (
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
	server.mux.Handle("")
}

func (Server *s) Run() {
	http.ListenAndServe(":"+port, s.mux)
}
