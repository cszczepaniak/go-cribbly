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

	g1, err := gameStore.Create(random.UUID(), random.UUID(), model.PrelimGame)
	require.NoError(t, err)

	g2, err := gameStore.Create(random.UUID(), random.UUID(), model.PrelimGame)
	require.NoError(t, err)

	g3, err := gameStore.Create(random.UUID(), random.UUID(), model.PrelimGame)
	require.NoError(t, err)

	g, err := gameStore.Get(g1.ID)
	require.NoError(t, err)
	assert.Equal(t, g1, g)

	g, err = gameStore.Get(g2.ID)
	require.NoError(t, err)
	assert.Equal(t, g2, g)

	g, err = gameStore.Get(g3.ID)
	require.NoError(t, err)
	assert.Equal(t, g3, g)

	gs, err := gameStore.GetAll()
	require.NoError(t, err)
	assert.Contains(t, gs, g1)
	assert.Contains(t, gs, g2)
	assert.Contains(t, gs, g3)
}
