package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/go-cribbly/internal/model"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/bytestore"
	"github.com/cszczepaniak/go-cribbly/internal/persistence/games"
	"github.com/cszczepaniak/go-cribbly/server/handlers"
)

func newTestServer(t *testing.T) *httptest.Server {
	byteStore := bytestore.NewMemoryByteStore()
	gameStore := games.NewS3GameStore(byteStore)
	handler := handlers.NewRequestHandler(gameStore)
	server := NewServer(handler)

	s := httptest.NewServer(server)
	t.Cleanup(s.Close)

	return s
}

func TestPing(t *testing.T) {
	s := newTestServer(t)

	resp, err := http.DefaultClient.Get(s.URL + `/ping`)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGames(t *testing.T) {
	s := newTestServer(t)

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
