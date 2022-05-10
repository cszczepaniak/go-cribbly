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
		ID: random.UUID(),
	}
	_, err := pcfg.GameStore.Create(g)
	require.NoError(t, err)
	tm := model.Team{
		ID: random.UUID(),
	}
	_, err = pcfg.TeamStore.Create(tm)
	require.NoError(t, err)

	var gr model.GameResult
	t.Run(`create`, func(t *testing.T) {
		r := strings.NewReader(fmt.Sprintf(`{
			"winner": %q,
			"loser_score": 111
		}`, tm.ID))
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
			"winner": %q,
			"loser_score": 111
		}`, tm.ID))
		resp, err := http.DefaultClient.Post(s.URL+`/games/`+random.UUID()+`/result`, ``, r)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run(`create for nonexistent team`, func(t *testing.T) {
		r := strings.NewReader(fmt.Sprintf(`{
			"winner": %q,
			"loser_score": 111
		}`, random.UUID()))
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

	t.Run(`get all`, func(t *testing.T) {
		resp, err := http.DefaultClient.Get(s.URL + `/games/` + g.ID + `/results`)
		require.NoError(t, err)
		defer resp.Body.Close()

		bs, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.NotEmpty(t, bs)

		var results []model.GameResult
		require.NoError(t, json.Unmarshal(bs, &results))
		require.Len(t, results, 1)
		assert.Equal(t, gr.ID, results[0].ID)
		assert.Equal(t, gr.Winner, results[0].Winner)
		assert.Equal(t, g.ID, results[0].GameID)
		assert.Equal(t, 111, results[0].LoserScore)
	})
}
