package server

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/go-cribbly/internal/model"
	"github.com/cszczepaniak/go-cribbly/internal/random"
)

func TestStandings(t *testing.T) {
	s, pcfg := newTestServer(t)

	_, err := pcfg.GameResultStore.Create(model.GameResult{
		ID:         random.UUID(),
		GameID:     random.UUID(),
		Winner:     random.UUID(),
		Loser:      random.UUID(),
		LoserScore: 100,
	})
	require.NoError(t, err)
	_, err = pcfg.GameResultStore.Create(model.GameResult{
		ID:         random.UUID(),
		GameID:     random.UUID(),
		Winner:     random.UUID(),
		Loser:      random.UUID(),
		LoserScore: 100,
	})
	require.NoError(t, err)
	_, err = pcfg.GameResultStore.Create(model.GameResult{
		ID:         random.UUID(),
		GameID:     random.UUID(),
		Winner:     random.UUID(),
		Loser:      random.UUID(),
		LoserScore: 100,
	})
	require.NoError(t, err)

	resp, err := http.DefaultClient.Get(s.URL + `/standings`)
	require.NoError(t, err)
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.NotEmpty(t, bs)

	var result []model.Standing
	require.NoError(t, json.Unmarshal(bs, &result))

	assert.Len(t, result, 6)
}
