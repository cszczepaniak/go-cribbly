package games

import (
	"github.com/cszczepaniak/go-cribbly/internal/model"
)

type GameStore interface {
	Create(teamAID, teamBID string, kind model.GameKind) (model.Game, error)
	Get(id string) (model.Game, error)
	GetAll() ([]model.Game, error)
}
