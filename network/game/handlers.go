package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cszczepaniak/go-cribbly/data/persistence"
)

type GameHandler struct {
	logger *log.Logger
	pcfg   *persistence.Config
}

func NewGameHandler(logger *log.Logger, pcfg *persistence.Config) *GameHandler {
	return &GameHandler{
		logger: logger,
		pcfg:   pcfg,
	}
}

func (gh *GameHandler) RegisterRoutes(r *mux.Router) {
	gh.logger.Println(`registering game routes...`)
	gamesRouter := r.Path(`/game`).Subrouter()
	gamesRouter.HandleFunc(``, gh.handleGetAll).Methods(http.MethodGet)
	gamesRouter.HandleFunc(``, gh.handleCreate).Methods(http.MethodPost)
}

func (gh *GameHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var teams []string
	err := json.NewDecoder(r.Body).Decode(&teams)
	if err != nil {
		gh.logger.Println(err)
		http.Error(w, `failed to unmarshal create game request`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(teams) != 2 {
		gh.logger.Println(`expected 2 teams, got`, len(teams))
		http.Error(w, `must provide exactly two teams`, http.StatusBadRequest)
		return
	}

	id, err := gh.pcfg.GameRepository.Create(teams[0], teams[1])
	if err != nil {
		gh.logger.Println(err)
		http.Error(w, `failed to create game`, http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, id)
}

func (gh *GameHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	games, err := gh.pcfg.GameRepository.GetAll()
	if err != nil {
		http.Error(w, `failed to get games`, http.StatusInternalServerError)
		return

	}

	err = json.NewEncoder(w).Encode(games)
	if err != nil {
		http.Error(w, `failed to marshal games`, http.StatusInternalServerError)
		return
	}
}
