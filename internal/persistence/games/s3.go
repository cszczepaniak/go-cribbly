package games

import (
	"bytes"
	"encoding/json"

	"github.com/cszczepaniak/go-cribbly/internal/model"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/bytestore"
	"github.com/cszczepaniak/go-cribbly/internal/random"
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

func (s *s3GameStore) Create(teamAID string, teamBID string, kind model.GameKind) (model.Game, error) {
	g := model.Game{
		ID: random.UUID(),
		TeamIDs: []string{
			teamAID, teamBID,
		},
		Kind: kind,
	}

	bs, err := json.Marshal(g)
	if err != nil {
		return model.Game{}, err
	}

	err = s.byteStore.Put(gameKey(g.ID), bytes.NewReader(bs))
	if err != nil {
		return model.Game{}, err
	}

	return g, nil
}

// Get implements GameStore
func (s *s3GameStore) Get(id string) (model.Game, error) {
	bs, err := s.byteStore.Get(gameKey(id))
	if err != nil {
		return model.Game{}, err
	}

	var g model.Game
	err = json.Unmarshal(bs, &g)
	if err != nil {
		return model.Game{}, err
	}
	return g, nil
}

// GetAll implements GameStore
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
