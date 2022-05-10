package handlers

import (
	"github.com/cszczepaniak/go-cribbly/internal/persistence"
	"github.com/cszczepaniak/go-cribbly/service/standings"
)

type RequestHandler struct {
	pcfg             *persistence.Config
	standingsService *standings.StandingsService
}

func NewRequestHandler(pcfg *persistence.Config, standingsService *standings.StandingsService) RequestHandler {
	return RequestHandler{
		pcfg:             pcfg,
		standingsService: standingsService,
	}
}
