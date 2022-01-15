package network

import (
	"github.com/cszczepaniak/go-cribbly/network/game"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	gh := game.NewGameHandler()
	gh.RegisterRoutes(router)

	return router
}
