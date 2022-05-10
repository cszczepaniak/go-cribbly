package gameresults

import (
	"github.com/cszczepaniak/go-cribbly/internal/model"
)

type GameResultStore interface {
	Create(model.GameResult) (model.GameResult, error)
	Get(id string) (model.GameResult, error)
	GetAll() ([]model.GameResult, error)
}
