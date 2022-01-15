package repository

import (
	"github.com/cszczepaniak/go-cribbly/data/game"
)

type GameRepository interface {
	Create(teamA string, teamB string) (string, error)
	Update(game *game.Game) error
	Delete(id string) error
	Get(id string) (*game.Game, error)
	GetAll() ([]*game.Game, error)
}
