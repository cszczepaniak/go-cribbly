package games

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/go-cribbly/internal/model"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/bytestore"
	"github.com/cszczepaniak/go-cribbly/internal/random"
)

func TestGames(t *testing.T) {
	byteStore := bytestore.NewMemoryByteStore()
	gameStore := NewS3GameStore(byteStore)

	e1, err := gameStore.Create(model.Game{TeamIDs: []string{random.UUID(), random.UUID()}, Kind: model.PrelimGame})
	require.NoError(t, err)

	e2, err := gameStore.Create(model.Game{TeamIDs: []string{random.UUID(), random.UUID()}, Kind: model.PrelimGame})
	require.NoError(t, err)

	e3, err := gameStore.Create(model.Game{TeamIDs: []string{random.UUID(), random.UUID()}, Kind: model.PrelimGame})
	require.NoError(t, err)

	e, err := gameStore.Get(e1.ID)
	require.NoError(t, err)
	assert.Equal(t, e1, e)

	e, err = gameStore.Get(e2.ID)
	require.NoError(t, err)
	assert.Equal(t, e2, e)

	e, err = gameStore.Get(e3.ID)
	require.NoError(t, err)
	assert.Equal(t, e3, e)

	es, err := gameStore.GetAll()
	require.NoError(t, err)
	assert.Contains(t, es, e1)
	assert.Contains(t, es, e2)
	assert.Contains(t, es, e3)
}
