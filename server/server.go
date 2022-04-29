package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"

	"github.com/cszczepaniak/go-cribbly/internal/persistence/bytestore"
)

type Config struct {
	IsLambda bool
}

type server struct {
	eng       *gin.Engine
	byteStore bytestore.ByteStore
}

func NewServer(bytestore bytestore.ByteStore) http.Handler {
	eng := gin.Default()
	eng.Use(cors.AllowAll())
	s := &server{
		eng:       eng,
		byteStore: bytestore,
	}

	s.registerRoutes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.eng.ServeHTTP(w, r)
}
