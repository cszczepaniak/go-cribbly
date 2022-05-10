package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/go-cribbly/internal/model"
)

func TestGames(t *testing.T) {
	s, _ := newTestServer(t)

	var g model.Game
	t.Run(`create`, func(t *testing.T) {
		r := strings.NewReader(`{
			"teams": ["abc", "def"],
			"kind": "prelim"
		}`)
		resp, err := http.DefaultClient.Post(s.URL+`/games`, ``, r)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		bs, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.NotEmpty(t, bs)

		require.NoError(t, json.Unmarshal(bs, &g))
	})

	t.Run(`get`, func(t *testing.T) {
		resp, err := http.DefaultClient.Get(s.URL + `/games/` + g.ID)
		require.NoError(t, err)
		defer resp.Body.Close()

		bs, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.NotEmpty(t, bs)

		var game model.Game
		require.NoError(t, json.Unmarshal(bs, &game))

		assert.Equal(t, g.ID, game.ID)
		assert.Equal(t, []string{`abc`, `def`}, game.TeamIDs)
		assert.Equal(t, model.GameKind(`prelim`), game.Kind)
	})
}
