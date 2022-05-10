package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/go-cribbly/internal/model"
	"github.com/cszczepaniak/go-cribbly/internal/random"
)

func TestGameResults(t *testing.T) {
	s, pcfg := newTestServer(t)

	g := model.Game{
		ID:      random.UUID(),
		TeamIDs: []string{random.UUID(), random.UUID()},
		Kind:    model.PrelimGame,
	}
	_, err := pcfg.GameStore.Create(g)
	require.NoError(t, err)

	var gr model.GameResult
	t.Run(`create`, func(t *testing.T) {
		r := strings.NewReader(fmt.Sprintf(`{
			"game_id": %q,
			"winner": %q,
			"loser_score": 111
		}`, g.ID, random.UUID()))
		resp, err := http.DefaultClient.Post(s.URL+`/games/`+g.ID+`/result`, ``, r)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		bs, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.NotEmpty(t, bs)

		require.NoError(t, json.Unmarshal(bs, &gr))
	})

	t.Run(`create for nonexistent game`, func(t *testing.T) {
		r := strings.NewReader(fmt.Sprintf(`{
			"game_id": %q,
			"winner": %q,
			"loser_score": 111
		}`, random.UUID(), random.UUID()))
		resp, err := http.DefaultClient.Post(s.URL+`/games/`+g.ID+`/result`, ``, r)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run(`get`, func(t *testing.T) {
		resp, err := http.DefaultClient.Get(s.URL + `/games/` + g.ID + `/result`)
		require.NoError(t, err)
		defer resp.Body.Close()

		bs, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.NotEmpty(t, bs)

		var result model.GameResult
		require.NoError(t, json.Unmarshal(bs, &result))

		assert.Equal(t, gr.ID, result.ID)
		assert.Equal(t, gr.Winner, result.Winner)
		assert.Equal(t, g.ID, result.GameID)
		assert.Equal(t, 111, result.LoserScore)
	})
}
