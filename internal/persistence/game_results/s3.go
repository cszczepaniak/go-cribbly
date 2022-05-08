package gameResults

import (
	"encoding/json"

	"github.com/cszczepaniak/go-cribbly/internal/model"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/bytestore"
	"github.com/cszczepaniak/go-cribbly/internal/random"
)

const gameResultsPrefix = `gameResults/`

func gameResultKey(id string) string {
	return gameResultsPrefix + id + `.json`
}

type s3GameResultStore struct {
	byteStore bytestore.ByteStore
}

var _ GameResultStore = (*s3GameResultStore)(nil)

func NewS3GameResultStore(byteStore bytestore.ByteStore) *s3GameResultStore {
	return &s3GameResultStore{
		byteStore: byteStore,
	}
}

func (s *s3GameResultStore) Create(e model.GameResult) (model.GameResult, error) {
	e.ID = random.UUID()
	err := s.byteStore.PutJSON(gameResultKey(e.ID), e)
	if err != nil {
		return model.GameResult{}, err
	}

	return e, nil
}

func (s *s3GameResultStore) Get(id string) (model.GameResult, error) {
	var v model.GameResult

	err := s.byteStore.GetJSON(gameResultKey(id), &v)
	if err != nil {
		return model.GameResult{}, err
	}

	return v, nil
}

func (s *s3GameResultStore) GetAll() ([]model.GameResult, error) {
	keyToPayload, err := s.byteStore.GetWithPrefix(gameResultsPrefix)
	if err != nil {
		return nil, err
	}

	res := make([]model.GameResult, len(keyToPayload))
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
