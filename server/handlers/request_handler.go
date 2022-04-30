package handlers

import "github.com/cszczepaniak/go-cribbly/internal/persistence/games"

type RequestHandler struct {
	gameStore games.GameStore
}

func NewRequestHandler(gameStore games.GameStore) RequestHandler {
	return RequestHandler{
		gameStore: gameStore,
	}
}
