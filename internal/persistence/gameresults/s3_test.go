package gameresults

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/go-cribbly/internal/model"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/bytestore"
	"github.com/cszczepaniak/go-cribbly/internal/random"
)

func TestGameResults(t *testing.T) {
	byteStore := bytestore.NewMemoryByteStore()
	gameResultStore := NewS3GameResultStore(byteStore)

	e1, err := gameResultStore.Create(model.GameResult{GameID: random.UUID(), Winner: random.UUID(), LoserScore: random.Int()})
	require.NoError(t, err)

	e2, err := gameResultStore.Create(model.GameResult{GameID: random.UUID(), Winner: random.UUID(), LoserScore: random.Int()})
	require.NoError(t, err)

	e3, err := gameResultStore.Create(model.GameResult{GameID: random.UUID(), Winner: random.UUID(), LoserScore: random.Int()})
	require.NoError(t, err)

	e, err := gameResultStore.Get(e1.ID)
	require.NoError(t, err)
	assert.Equal(t, e1, e)

	e, err = gameResultStore.Get(e2.ID)
	require.NoError(t, err)
	assert.Equal(t, e2, e)

	e, err = gameResultStore.Get(e3.ID)
	require.NoError(t, err)
	assert.Equal(t, e3, e)

	es, err := gameResultStore.GetAll()
	require.NoError(t, err)
	assert.Contains(t, es, e1)
	assert.Contains(t, es, e2)
	assert.Contains(t, es, e3)
}
