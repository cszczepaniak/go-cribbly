package repository

import "github.com/cszczepaniak/go-cribbly/game/model"

type GameRepository interface {
	Create(teamA string, teamB string) (string, error)
	Update(g *model.Game) error
	Delete(id string) error
	Get(id string) (*model.Game, error)
	GetAll() ([]*model.Game, error)
}
