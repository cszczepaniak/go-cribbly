//go:build !prod

package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cszczepaniak/go-cribbly/server/handlers"
)

func NewTestServer(handler handlers.RequestHandler) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	eng := gin.New()
	return newServer(handler, eng)
}
