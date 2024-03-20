package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	Port        string
	mux         *http.ServeMux
	middlewares []http.HandlerFunc
}

func NewServer(port string) *Server {
	mux := http.NewServeMux()
	return &Server{Port: port, mux: mux}
}

func (s *Server) AddRoute(method string, path string, handler http.HandlerFunc) {
	s.mux.HandleFunc(method+" "+path, func(w http.ResponseWriter, r *http.Request) {
		for _, middleware := range s.middlewares {
			middleware(w, r)
		}
		handler(w, r)
	})
}

func (s *Server) AddHandler(method string, path string, handler http.Handler) {
	s.mux.Handle(method+" "+path, handler)
}

func (s *Server) Start() {
	fmt.Println("Server running on port", s.Port)
	err := http.ListenAndServe(":"+s.Port, s.mux)
	if err != nil {
		panic(err)
	}
}
