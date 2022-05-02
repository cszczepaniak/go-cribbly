package teams

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/go-cribbly/internal/model"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/bytestore"
	"github.com/cszczepaniak/go-cribbly/internal/random"
)

func TestTeams(t *testing.T) {
	byteStore := bytestore.NewMemoryByteStore()
	teamStore := NewS3TeamStore(byteStore)

	ps := make([]model.Player, 0, 6)
	for i := 0; i < 6; i++ {
		ps = append(ps, model.Player{
			FirstName: random.UUID(),
			LastName:  random.UUID(),
		})
	}

	t1, err := teamStore.Create(ps[0], ps[1])
	require.NoError(t, err)

	t2, err := teamStore.Create(ps[2], ps[3])
	require.NoError(t, err)

	t3, err := teamStore.Create(ps[4], ps[5])
	require.NoError(t, err)

	team, err := teamStore.Get(t1.ID)
	require.NoError(t, err)
	assert.Equal(t, t1, team)

	team, err = teamStore.Get(t2.ID)
	require.NoError(t, err)
	assert.Equal(t, t2, team)

	team, err = teamStore.Get(t3.ID)
	require.NoError(t, err)
	assert.Equal(t, t3, team)

	teams, err := teamStore.GetAll()
	require.NoError(t, err)
	assert.Contains(t, teams, t1)
	assert.Contains(t, teams, t2)
	assert.Contains(t, teams, t3)
}
