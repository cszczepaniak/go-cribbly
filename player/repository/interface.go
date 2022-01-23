package repository

import "github.com/cszczepaniak/go-cribbly/player/model"

type PlayerRepository interface {
	Create( /* TODO */ ) (string, error)
	Update(p *model.Player) error
	Delete(id string) error
	Get(id string) (*model.Player, error)
	GetAll() ([]*model.Player, error)
}
