package games

import (
	"encoding/json"

	"github.com/cszczepaniak/go-cribbly/internal/model"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/bytestore"
)

const gamesPrefix = `games/`

func gameKey(id string) string {
	return gamesPrefix + id + `.json`
}

type s3GameStore struct {
	byteStore bytestore.ByteStore
}

var _ GameStore = (*s3GameStore)(nil)

func NewS3GameStore(byteStore bytestore.ByteStore) *s3GameStore {
	return &s3GameStore{
		byteStore: byteStore,
	}
}

func (s *s3GameStore) Create(e model.Game) (model.Game, error) {
	err := s.byteStore.PutJSON(gameKey(e.ID), e)
	if err != nil {
		return model.Game{}, err
	}

	return e, nil
}

func (s *s3GameStore) Get(id string) (model.Game, error) {
	var v model.Game

	err := s.byteStore.GetJSON(gameKey(id), &v)
	if err != nil {
		return model.Game{}, err
	}

	return v, nil
}

func (s *s3GameStore) GetAll() ([]model.Game, error) {
	keyToPayload, err := s.byteStore.GetWithPrefix(gamesPrefix)
	if err != nil {
		return nil, err
	}

	res := make([]model.Game, len(keyToPayload))
	i := 0
	for _, p := range keyToPayload {
		err := json.Unmarshal(p, &res[i])
		if err != nil {
			return nil, err
		}
		i++
	}

	return res, nil
}
