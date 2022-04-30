package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"

	"github.com/cszczepaniak/go-cribbly/server/handlers"
)

type Config struct {
	IsLambda bool
}

type server struct {
	eng            *gin.Engine
	requestHandler handlers.RequestHandler
}

func NewServer(handler handlers.RequestHandler) http.Handler {
	eng := gin.Default()
	eng.Use(cors.AllowAll())
	s := &server{
		eng:            eng,
		requestHandler: handler,
	}

	s.registerRoutes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.eng.ServeHTTP(w, r)
}
