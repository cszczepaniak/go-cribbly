package game

import (
	"encoding/json"
	"net/http"

	"github.com/cszczepaniak/go-cribbly/data/game"
	"github.com/gorilla/mux"
)

type GameHandler struct{}

func NewGameHandler() *GameHandler {
	return &GameHandler{}
}

func (gh *GameHandler) RegisterRoutes(r *mux.Router) {
	gameRoute := r.Path(`/game`)
	gameRoute.Methods(http.MethodGet).HandlerFunc(gh.handleGetAll)
}

func (gh *GameHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	games := []game.Game{{
		Scores: map[string]int{`a`: 93, `b`: 121},
	}}
	err := json.NewEncoder(w).Encode(games)
	if err != nil {
		http.Error(w, `failed to marshal games`, http.StatusInternalServerError)
		return
	}
}
