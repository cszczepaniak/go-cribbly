package server

import (
	"fmt"
	"net/http"
)

func (s *Server) RegisterRoutes() {
	s.router.Get(`/ping`, func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `pong`)
	})
}
