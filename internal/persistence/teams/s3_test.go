package teams

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/go-cribbly/internal/model"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/bytestore"
	"github.com/cszczepaniak/go-cribbly/internal/random"
)

func randomPlayer() model.Player {
	return model.Player{
		FirstName: random.UUID(),
		LastName:  random.UUID(),
	}
}

func TestTeams(t *testing.T) {
	byteStore := bytestore.NewMemoryByteStore()
	teamStore := NewS3TeamStore(byteStore)

	e1, err := teamStore.Create(randomPlayer(), randomPlayer())
	require.NoError(t, err)

	e2, err := teamStore.Create(randomPlayer(), randomPlayer())
	require.NoError(t, err)

	e3, err := teamStore.Create(randomPlayer(), randomPlayer())
	require.NoError(t, err)

	e, err := teamStore.Get(e1.ID)
	require.NoError(t, err)
	assert.Equal(t, e1, e)

	e, err = teamStore.Get(e2.ID)
	require.NoError(t, err)
	assert.Equal(t, e2, e)

	e, err = teamStore.Get(e3.ID)
	require.NoError(t, err)
	assert.Equal(t, e3, e)

	es, err := teamStore.GetAll()
	require.NoError(t, err)
	assert.Contains(t, es, e1)
	assert.Contains(t, es, e2)
	assert.Contains(t, es, e3)
}
