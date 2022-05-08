package games

import (
	"github.com/cszczepaniak/go-cribbly/internal/model"
)

type GameStore interface {
	Create(teamA, teamB string, kind model.GameKind) (model.Game, error)
	Get(id string) (model.Game, error)
	GetAll() ([]model.Game, error)
}
