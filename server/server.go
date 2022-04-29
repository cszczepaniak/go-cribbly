package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"

	"github.com/cszczepaniak/go-cribbly/internal/persistence/s3"
)

type Config struct {
	IsLambda bool
}

type server struct {
	eng      *gin.Engine
	s3Client s3.ByteStore
}

func NewServer(s3Client s3.ByteStore) http.Handler {
	eng := gin.Default()
	eng.Use(cors.AllowAll())
	s := &server{
		eng:      eng,
		s3Client: s3Client,
	}

	s.registerRoutes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.eng.ServeHTTP(w, r)
}
