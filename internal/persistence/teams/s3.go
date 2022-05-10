package teams

import (
	"encoding/json"

	"github.com/cszczepaniak/go-cribbly/internal/model"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/bytestore"
)

const teamsPrefix = `teams/`

func teamKey(id string) string {
	return teamsPrefix + id + `.json`
}

type s3TeamStore struct {
	byteStore bytestore.ByteStore
}

var _ TeamStore = (*s3TeamStore)(nil)

func NewS3TeamStore(byteStore bytestore.ByteStore) *s3TeamStore {
	return &s3TeamStore{
		byteStore: byteStore,
	}
}

func (s *s3TeamStore) Create(e model.Team) (model.Team, error) {
	err := s.byteStore.PutJSON(teamKey(e.ID), e)
	if err != nil {
		return model.Team{}, err
	}

	return e, nil
}

func (s *s3TeamStore) Get(id string) (model.Team, error) {
	var v model.Team

	err := s.byteStore.GetJSON(teamKey(id), &v)
	if err != nil {
		return model.Team{}, err
	}

	return v, nil
}

func (s *s3TeamStore) GetAll() ([]model.Team, error) {
	keyToPayload, err := s.byteStore.GetWithPrefix(teamsPrefix)
	if err != nil {
		return nil, err
	}

	res := make([]model.Team, len(keyToPayload))
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
