package server

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) AddMiddleware(middleware http.HandlerFunc) {
	s.middlewares = append(s.middlewares, middleware)
}

var RequestIdMiddleware http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Request-ID", uuid.New().String())
}
