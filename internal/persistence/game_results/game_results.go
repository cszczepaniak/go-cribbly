package gameResults

import (
	"github.com/cszczepaniak/go-cribbly/internal/model"
)

type GameResultStore interface {
	Create(gameID, winner string, scoreDifference int) (model.GameResult, error)
	Get(id string) (model.GameResult, error)
	GetAll() ([]model.GameResult, error)
}
