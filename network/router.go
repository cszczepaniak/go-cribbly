package network

import (
	"log"

	"github.com/gorilla/mux"

	"github.com/cszczepaniak/go-cribbly/data/persistence"
	"github.com/cszczepaniak/go-cribbly/network/game"
)

func SetupRouter(logger *log.Logger, pcfg *persistence.Config) *mux.Router {
	router := mux.NewRouter()

	gh := game.NewGameHandler(logger, pcfg)
	gh.RegisterRoutes(router)

	return router
}
