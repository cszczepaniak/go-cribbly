package repository

import (
	"errors"
	"sync"

	"github.com/cszczepaniak/go-cribbly/common/random"
	"github.com/cszczepaniak/go-cribbly/player/model"
)

var errPlayerNotFound = errors.New(`player not found`)

type memoryPlayerRepository struct {
	lock    sync.RWMutex
	players map[string]*model.Player
}

func NewMemory() PlayerRepository {
	return &memoryPlayerRepository{
		players: make(map[string]*model.Player),
	}
}

func (r *memoryPlayerRepository) Create( /* TODO */ ) (string, error) {
	id := random.UUID()
	p := model.NewPlayer(id)

	r.writePlayer(p)
	return id, nil
}

func (r *memoryPlayerRepository) Delete(id string) error {
	if _, ok := r.getPlayer(id); !ok {
		return errPlayerNotFound
	}

	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.players, id)

	return nil
}

func (r *memoryPlayerRepository) Get(id string) (*model.Player, error) {
	p, ok := r.getPlayer(id)
	if !ok {
		return nil, errPlayerNotFound
	}
	return p, nil
}

func (r *memoryPlayerRepository) GetAll() ([]*model.Player, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	players := make([]*model.Player, 0, len(r.players))
	for _, p := range r.players {
		players = append(players, p)
	}
	return players, nil
}

func (r *memoryPlayerRepository) Update(p *model.Player) error {
	r.lock.RLock()
	_, ok := r.players[p.ID]
	r.lock.RUnlock()

	if !ok {
		return errPlayerNotFound
	}

	r.writePlayer(p)
	return nil
}

func (r *memoryPlayerRepository) writePlayer(p *model.Player) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.players[p.ID] = p
}

func (r *memoryPlayerRepository) getPlayer(id string) (*model.Player, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	p, ok := r.players[id]
	return p, ok
}
