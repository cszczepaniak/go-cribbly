package network

import (
	"log"

	"github.com/gorilla/mux"

	game "github.com/cszczepaniak/go-cribbly/game/network"
	"github.com/cszczepaniak/go-cribbly/persistence"
)

func SetupRouter(logger *log.Logger, pcfg *persistence.Config) *mux.Router {
	router := mux.NewRouter()

	gh := game.NewGameHandler(logger, pcfg)
	gh.RegisterRoutes(router)

	return router
}
