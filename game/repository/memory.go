package repository

import (
	"errors"
	"sync"

	"github.com/cszczepaniak/go-cribbly/common/random"
	"github.com/cszczepaniak/go-cribbly/game/model"
)

var errGameNotFound = errors.New(`game not found`)

type memoryGameRepository struct {
	lock  sync.RWMutex
	games map[string]*model.Game
}

func NewMemory() GameRepository {
	return &memoryGameRepository{
		games: make(map[string]*model.Game),
	}
}

func (r *memoryGameRepository) Create(teamA, teamB string) (string, error) {
	id := random.UUID()
	g := model.NewGame(id, teamA, teamB)

	r.writeGame(g)
	return id, nil
}

func (r *memoryGameRepository) Delete(id string) error {
	if _, ok := r.getGame(id); !ok {
		return errGameNotFound
	}

	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.games, id)

	return nil
}

func (r *memoryGameRepository) Get(id string) (*model.Game, error) {
	g, ok := r.getGame(id)
	if !ok {
		return nil, errGameNotFound
	}
	return g, nil
}

func (r *memoryGameRepository) GetAll() ([]*model.Game, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	games := make([]*model.Game, 0, len(r.games))
	for _, g := range r.games {
		games = append(games, g)
	}
	return games, nil
}

func (r *memoryGameRepository) Update(g *model.Game) error {
	r.lock.RLock()
	_, ok := r.games[g.ID]
	r.lock.RUnlock()

	if !ok {
		return errGameNotFound
	}

	r.writeGame(g)
	return nil
}

func (r *memoryGameRepository) writeGame(g *model.Game) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.games[g.ID] = g
}

func (r *memoryGameRepository) getGame(id string) (*model.Game, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	g, ok := r.games[id]
	return g, ok
}
