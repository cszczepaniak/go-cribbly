package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	router *chi.Mux
	logger *log.Logger
}

func NewServer(l *log.Logger) *Server {
	return &Server{
		router: chi.NewRouter(),
		logger: l,
	}
}

func (s *Server) Serve() error {
	s.router.Use(middleware.Logger, middleware.RequestID)
	s.RegisterRoutes()

	return http.ListenAndServe(`:8080`, s.router)
}
