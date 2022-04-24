package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

type Config struct {
	IsLambda bool
}

type server struct {
	eng *gin.Engine
}

func NewServer() http.Handler {
	eng := gin.Default()
	eng.Use(cors.AllowAll())
	s := &server{
		eng: eng,
	}

	s.registerRoutes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.eng.ServeHTTP(w, r)
}
