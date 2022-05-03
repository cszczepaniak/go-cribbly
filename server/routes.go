package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *server) registerRoutes() {
	s.eng.GET(`/ping`, func(ctx *gin.Context) {
		ctx.String(http.StatusOK, `pong`)
	})

	games := s.eng.Group(`/games`)
	games.GET(``, s.requestHandler.HandleGetAllGames)
	games.GET(`/:id`, s.requestHandler.HandleGetGame)
	games.POST(``, s.requestHandler.HandleCreateGame)

	teams := s.eng.Group(`/teams`)
	teams.GET(``, s.requestHandler.HandleGetAllTeams)
	teams.GET(`/:id`, s.requestHandler.HandleGetTeam)
	teams.POST(``, s.requestHandler.HandleCreateTeam)
}
