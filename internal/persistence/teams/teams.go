package teams

import (
	"github.com/cszczepaniak/go-cribbly/internal/model"
)

type TeamStore interface {
	Create(playerA, playerB model.Player) (model.Team, error)
	Get(id string) (model.Team, error)
	GetAll() ([]model.Team, error)
}
