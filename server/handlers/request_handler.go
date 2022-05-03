package handlers

import (
	"github.com/cszczepaniak/go-cribbly/internal/persistence"
)

type RequestHandler struct {
	pcfg *persistence.Config
}

func NewRequestHandler(pcfg *persistence.Config) RequestHandler {
	return RequestHandler{
		pcfg: pcfg,
	}
}
