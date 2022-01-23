package network

import (
	"log"

	"github.com/gorilla/mux"

	"github.com/cszczepaniak/go-cribbly/persistence"
)

type PlayerHandler struct {
	logger *log.Logger
	pcfg   *persistence.Config
}

func NewPlayerHandler(logger *log.Logger, pcfg *persistence.Config) *PlayerHandler {
	return &PlayerHandler{
		logger: logger,
		pcfg:   pcfg,
	}
}

func (ph *PlayerHandler) RegisterRoutes(r *mux.Router) {
	playersRouter := r.Path(`/players`).Subrouter()

	// TODO add players routes
}
