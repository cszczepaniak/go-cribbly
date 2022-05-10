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

func TestTeams(t *testing.T) {
	s := newTestServer(t)

	var team model.Team
	t.Run(`create`, func(t *testing.T) {
		r := strings.NewReader(`{
			"players": [
				{"first_name": "foo", "last_name": "bar"},
				{"first_name": "baz", "last_name": "qux"}
			]
		}`)

		resp, err := http.DefaultClient.Post(s.URL+`/teams`, ``, r)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		bs, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.NotEmpty(t, bs)

		require.NoError(t, json.Unmarshal(bs, &team))
	})

	t.Run(`get`, func(t *testing.T) {
		resp, err := http.DefaultClient.Get(s.URL + `/teams/` + team.ID)
		require.NoError(t, err)
		defer resp.Body.Close()

		bs, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.NotEmpty(t, bs)

		var actualTeam model.Team
		require.NoError(t, json.Unmarshal(bs, &actualTeam))

		assert.Equal(t, team.ID, actualTeam.ID)
		assert.Len(t, actualTeam.Players, 2)
		assert.Contains(t, actualTeam.Players, model.Player{FirstName: `foo`, LastName: `bar`})
		assert.Contains(t, actualTeam.Players, model.Player{FirstName: `baz`, LastName: `qux`})
	})
}
