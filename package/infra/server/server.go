package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	Port string
	mux  *http.ServeMux
}

func NewServer(port string) *Server {
	mux := http.NewServeMux()
	return &Server{Port: port, mux: mux}
}

func (s *Server) AddRoute(method string, path string, handler http.HandlerFunc) {
	s.mux.HandleFunc(method+" "+path, handler)
}

func (s *Server) Start() {
	fmt.Println("Server running on port", s.Port)
	err := http.ListenAndServe(":"+s.Port, s.mux)
	if err != nil {
		panic(err)
	}
}
